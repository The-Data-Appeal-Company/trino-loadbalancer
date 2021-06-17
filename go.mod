module github.com/The-Data-Appeal-Company/trino-loadbalancer

go 1.15

require (
	github.com/aws/aws-sdk-go v1.34.30
	github.com/docker/go-connections v0.4.0
	github.com/go-redis/redis/v8 v8.2.2
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.6.2
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/lib/pq v1.8.0
	github.com/montanaflynn/stats v0.6.6
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/testcontainers/testcontainers-go v0.9.0
	github.com/trinodb/trino-go-client v0.300.0
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/kubectl v0.21.1
)
