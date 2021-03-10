package discovery

import (
	"context"
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/emr"
	"net/url"
	"strings"
)

const (
	TrinoEmrDefaultPort     = 8889
	TrinoEmrDefaultProtocol = "http"
)

type AwsCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

type ClusterProvider struct {
	emrClient  *emr.EMR
	ec2Client  *ec2.EC2
	SelectTags map[string]string
}

func AwsEmrDiscovery(cred AwsCredentials) *ClusterProvider {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      &cred.Region,
		Credentials: credentials.NewStaticCredentials(cred.AccessKeyID, cred.SecretAccessKey, ""),
	}))

	return &ClusterProvider{
		emrClient: emr.New(sess),
		ec2Client: ec2.New(sess),
	}
}

func (c *ClusterProvider) Discover() ([]models.Coordinator, error) {
	masters, err := c.listTargetMasters(context.Background())
	if err != nil {
		return nil, err
	}

	filtered := make([]models.Coordinator, 0)
	for _, m := range masters {
		if containsAll(m.Tags, c.SelectTags) {
			filtered = append(filtered, m)
		}
	}

	return filtered, nil
}

func (c *ClusterProvider) listTargetMasters(ctx context.Context) ([]models.Coordinator, error) {

	coordinators := make([]models.Coordinator, 0)

	clusters, err := c.listTargetClusters(ctx)

	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters {

		master, err := c.getClusterMasterInstance(cluster)
		if err != nil {
			return nil, err
		}

		masterUrl, err := url.Parse(fmt.Sprintf("%s://%s:%d", TrinoEmrDefaultProtocol, master, TrinoEmrDefaultPort))
		if err != nil {
			return nil, err
		}
		coordinators = append(coordinators, models.Coordinator{
			Name:    *cluster.Cluster.Id,
			URL:     masterUrl,
			Tags:    tagsToMap(cluster.Cluster.Tags),
			Enabled: true,
		})
	}

	return coordinators, nil
}

func tagsToMap(tags []*emr.Tag) map[string]string {
	res := make(map[string]string)

	for _, t := range tags {
		res[*t.Key] = *t.Value
	}

	return res
}

func (c *ClusterProvider) listTargetClusters(ctx context.Context) ([]*emr.DescribeClusterOutput, error) {
	req := &emr.ListClustersInput{
		ClusterStates: aws.StringSlice([]string{"WAITING"}),
	}

	clusters := make([]*emr.DescribeClusterOutput, 0)
	err := c.emrClient.ListClustersPagesWithContext(ctx, req, func(output *emr.ListClustersOutput, b bool) bool {

		for _, cluster := range output.Clusters {

			if *cluster.Status.State != emr.ClusterStateWaiting && *cluster.Status.State != emr.ClusterStateRunning {
				continue
			}

			descr, _ := c.emrClient.DescribeCluster(&emr.DescribeClusterInput{
				ClusterId: cluster.Id,
			})

			if !hasTrinoInstallation(descr) {
				continue
			}

			clusters = append(clusters, descr)

		}
		return true
	})

	return clusters, err
}

func (c *ClusterProvider) getClusterMasterInstance(cluster *emr.DescribeClusterOutput) (string, error) {

	instanceCollectionType := cluster.Cluster.InstanceCollectionType

	if *instanceCollectionType == emr.InstanceCollectionTypeInstanceGroup {
		return c.getMasterInstanceForNodeGroup(cluster)
	} else if *instanceCollectionType == emr.InstanceCollectionTypeInstanceFleet {
		return c.getMasterInstanceForFleet(cluster)
	}

	return "", fmt.Errorf("unrecognized instance type %s", *instanceCollectionType)
}

func (c *ClusterProvider) getMasterInstanceForFleet(cluster *emr.DescribeClusterOutput) (string, error) {

	instances, err := c.emrClient.ListInstances(&emr.ListInstancesInput{
		ClusterId:         cluster.Cluster.Id,
		InstanceFleetType: aws.String(emr.InstanceFleetTypeMaster),
	})

	if err != nil {
		return "", err
	}

	if len(instances.Instances) == 0 {
		return "", fmt.Errorf("no master instance found for cluster %s", *cluster.Cluster.Id)
	}

	return *instances.Instances[0].PrivateIpAddress, nil
}

func (c *ClusterProvider) getMasterInstanceForNodeGroup(cluster *emr.DescribeClusterOutput) (string, error) {

	instanceGroups, err := c.emrClient.ListInstances(&emr.ListInstancesInput{
		ClusterId:          cluster.Cluster.Id,
		InstanceGroupTypes: []*string{aws.String(emr.InstanceGroupTypeMaster)},
	})

	if err != nil {
		return "", err
	}

	for _, group := range instanceGroups.Instances {

		instances, err := c.emrClient.ListInstances(&emr.ListInstancesInput{
			ClusterId:       cluster.Cluster.Id,
			InstanceGroupId: group.Id,
		})

		if err != nil {
			return "", err
		}

		if len(instances.Instances) == 0 {
			return "", fmt.Errorf("no master instances found for cluster %s", *cluster.Cluster.Id)
		}

		return *instances.Instances[0].PrivateIpAddress, nil
	}

	return "", fmt.Errorf("no master instance found for cluster %s", *cluster.Cluster.Id)
}

func hasTrinoInstallation(descr *emr.DescribeClusterOutput) bool {
	for _, application := range descr.Cluster.Applications {
		if strings.Contains(strings.ToLower(*application.Name), "trino") {
			return true
		}
	}
	return false
}
