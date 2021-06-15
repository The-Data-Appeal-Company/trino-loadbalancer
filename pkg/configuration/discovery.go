package configuration

import (
	"database/sql"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery/kubernetes"
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
	Provider string                    `json:"provider" yaml:"provider" mapstructure:"provider"`
	Enabled  bool                      `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Aws      AwsDiscoveryConfiguration `json:"aws" yaml:"aws" mapstructure:"aws"`
	K8s      K8sConfiguration          `json:"k8s" yaml:"k8s" mapstructure:"k8s"`
	Static   StaticConfiguration       `json:"static" yaml:"static" mapstructure:"static"`
}

type AwsDiscoveryConfiguration struct {
	AwsAccessKeyID string `json:"access_key_id" yaml:"access_key_id" mapstructure:"access_key_id"`
	AwsSecretKey   string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
	AwsRegion      string `json:"region" yaml:"region" mapstructure:"region"`
}

type K8sConfiguration struct {
	KubeConfig    string            `json:"kube_config" yaml:"kube_config" mapstructure:"kube_config"`
	ClusterDomain string            `json:"cluster_domain" yaml:"cluster_domain" mapstructure:"cluster_domain"`
	SelectorTags  map[string]string `json:"selector_tags" yaml:"selector_tags" mapstructure:"selector_tags"`
}

type StaticConfiguration struct {
	Clusters []struct {
		Name    string            `json:"name" yaml:"name" mapstructure:"name"`
		Url     string            `json:"url" yaml:"url" mapstructure:"url"`
		Tags    map[string]string `json:"tags" yaml:"tags" mapstructure:"tags"`
		Enabled bool              `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	} `json:"clusters" yaml:"clusters" mapstructure:"clusters"`
}

const (
	DiscoveryTypeAws = "aws-emr"
	DiscoveryK8s     = "k8s"
	DiscoveryStatic  = "static"
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

	if conf.Provider == DiscoveryStatic {
		coordinators := make([]models.Coordinator, len(conf.Static.Clusters))
		for i, c := range conf.Static.Clusters {
			uri, err := url.Parse(c.Url)
			if err != nil {
				return nil, err
			}
			coordinators[i] = models.Coordinator{
				Name:    c.Name,
				URL:     uri,
				Tags:    c.Tags,
				Enabled: c.Enabled,
			}
		}

		return discovery.NewStatic(coordinators...), nil
	}

	return nil, fmt.Errorf("no discovery for type: %s", conf.Provider)
}
