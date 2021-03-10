package factory

import (
	"database/sql"
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/discovery/kubernetes"
	_ "github.com/lib/pq"
	"net/url"
)

type DiscoveryStorageConfiguration struct {
	Db       string
	Host     string
	Port     int
	User     string
	SslMode  string
	Password string
}

func CreateDiscoveryStorage(conf DiscoveryStorageConfiguration) (discovery.Storage, error) {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", conf.User, url.QueryEscape(conf.Password), conf.Host, conf.Port, conf.Db, conf.SslMode)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return discovery.NewDatabaseStorage(db, discovery.DefaultDatabaseTableName), nil
}

type DiscoveryConfiguration struct {
	Provider string                    `json:"provider"`
	Enabled  bool                      `json:"enabled"`
	Aws      AwsDiscoveryConfiguration `json:"aws"`
	K8s      K8sConfiguration          `json:"k8s"`
}

type AwsDiscoveryConfiguration struct {
	AwsAccessKeyID string `json:"access_key_id"`
	AwsSecretKey   string `json:"secret_key"`
	AwsRegion      string `json:"region"`
}

type K8sConfiguration struct {
	KubeConfig    string            `json:"kube_config"`
	ClusterDomain string            `json:"cluster_domain"`
	SelectorTags  map[string]string `json:"selector_tags"`
}

const (
	DiscoveryTypeAws = "aws-emr"
	DiscoveryK8s     = "k8s"
)

func CreateCrossProviderDiscovery(configs []DiscoveryConfiguration) (discovery.Discovery, error) {

	discoveryProviders := make([]discovery.Discovery, 0)

	for _, config := range configs {
		if config.Enabled {
			discoveryProvider, err := CreateDiscoveryProvider(config)

			if err != nil {
				return nil, err
			}
			discoveryProviders = append(discoveryProviders, discoveryProvider)
		}
	}

	providerDiscovery := discovery.NewCrossProviderDiscovery(discoveryProviders)

	return providerDiscovery, nil
}

func CreateDiscoveryProvider(conf DiscoveryConfiguration) (discovery.Discovery, error) {
	if !conf.Enabled {
		return discovery.Noop(), nil
	}

	if conf.Provider == DiscoveryTypeAws {
		return discovery.AwsEmrDiscovery(discovery.AwsCredentials{
			AccessKeyID:     conf.Aws.AwsAccessKeyID,
			SecretAccessKey: conf.Aws.AwsSecretKey,
			Region:          conf.Aws.AwsRegion,
		}), nil
	}

	if conf.Provider == DiscoveryK8s {
		client, err := kubernetes.NewClient(&conf.K8s.KubeConfig)

		if err != nil {
			return nil, err
		}

		return discovery.NewK8sClusterProvider(client, conf.K8s.SelectorTags, conf.K8s.ClusterDomain), nil
	}

	return nil, fmt.Errorf("no discovery for type: %s", conf.Provider)
}
