module github.com/The-Data-Appeal-Company/trino-loadbalancer

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.30
	github.com/docker/go-connections v0.4.0
	github.com/go-redis/redis/v8 v8.2.2
	github.com/gorilla/mux v1.6.2
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/lib/pq v1.8.0
	github.com/prestodb/presto-go-client v0.0.0-20200302111820-5ec09431be26
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.6.1
	github.com/testcontainers/testcontainers-go v0.9.0
	github.com/trinodb/trino-go-client v0.300.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.19.3
	k8s.io/utils v0.0.0-20201027101359-01387209bb0d // indirect
)
