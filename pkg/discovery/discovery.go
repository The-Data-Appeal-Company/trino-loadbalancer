package discovery

import "github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"

type Discovery interface {
	Discover() ([]models.Coordinator, error)
}
