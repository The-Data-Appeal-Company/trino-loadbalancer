package factory

import (
	"database/sql"
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/discovery"
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
	Enabled bool
	Type    string
	Aws     AwsDiscoveryConfiguration
}

type AwsDiscoveryConfiguration struct {
	AwsAccessKeyID string
	AwsSecretKey   string
	AwsRegion      string
}

const (
	DiscoveryTypeAws = "aws-emr"
)

func CreateDiscovery(conf DiscoveryConfiguration) (discovery.Discovery, error) {
	if !conf.Enabled {
		return discovery.Noop(),nil
	}

	if conf.Type == DiscoveryTypeAws {
		return discovery.AwsEmrDiscovery(discovery.AwsCredentials{
			AccessKeyID:     conf.Aws.AwsAccessKeyID,
			SecretAccessKey: conf.Aws.AwsSecretKey,
			Region:          conf.Aws.AwsRegion,
		}), nil
	}

	return nil, fmt.Errorf("no discovery for type: %s", conf.Type)
}
