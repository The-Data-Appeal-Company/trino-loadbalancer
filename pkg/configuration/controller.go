package configuration

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/controller/components"
	"github.com/go-redis/redis/v8"
	"time"
)

type ControllerConf struct {
	Controller struct {
		Features struct {
			SlowWorkerDrainer SlowWorkerDrainerConf `json:"slow_worker_drainer" yaml:"slow_worker_drainer" mapstructure:"slow_worker_drainer"`
		} `json:"features" yaml:"features" mapstructure:"features"`
	} `json:"controller" yaml:"controller" mapstructure:"controller"`
}

type SlowWorkerDrainerConf struct {
	Enabled            bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	GracePeriodSeconds int    `json:"GracePeriodSeconds" yaml:"GracePeriodSeconds" mapstructure:"GracePeriodSeconds"`
	DryRun             bool   `json:"dryRun" yaml:"dryRun" mapstructure:"dryRun"`
	DrainThreshold     int64  `json:"drainThreshold" yaml:"drainThreshold" mapstructure:"drainThreshold"`
	Provider           string `json:"provider" yaml:"provider" mapstructure:"provider"`
	K8sProvider        struct {
		KubeConfig        string            `json:"config" yaml:"config" mapstructure:"config"`
		NamespaceSelector map[string]string `json:"namespaceSelector" yaml:"namespaceSelector" mapstructure:"namespaceSelector"`
	} `json:"k8s" yaml:"k8s" mapstructure:"k8s"`
}

func CreateHandlers(redisClient redis.UniversalClient, logger logging.Logger, conf ControllerConf) (components.QueryHandler, error) {
	handlers := make([]components.QueryHandler, 0)

	slowNodeDrainerConf := conf.Controller.Features.SlowWorkerDrainer
	if slowNodeDrainerConf.Enabled {
		slowWorkerHandler, err := createDrainSlowWorkerNodeHandler(redisClient, logger, slowNodeDrainerConf)
		if err != nil {
			return nil, err
		}

		handlers = append(handlers, slowWorkerHandler)
	}

	return components.NewMultiQueryComponent(handlers...), nil
}

func createDrainSlowWorkerNodeHandler(redisClient redis.UniversalClient, logger logging.Logger, slowNodeDrainerConf SlowWorkerDrainerConf) (*components.SlowNodeDrainer, error) {
	analyzer := components.NewTrinoSlowNodeAnalyzer()
	nodeDrainer, err := createNodeDrainer(slowNodeDrainerConf, logger)
	if err != nil {
		return nil, err
	}
	slowNodeMarker := components.NewRedisSlowNodeMarker(redisClient)
	return components.NewSlowNodeDrainer(analyzer, nodeDrainer, slowNodeMarker, components.SlowNodeDrainerConf{DrainThreshold: slowNodeDrainerConf.DrainThreshold}, logger), nil
}

const ControllerK8sProvider = "k8s"

func createNodeDrainer(conf SlowWorkerDrainerConf, logger logging.Logger) (components.NodeDrainer, error) {
	switch conf.Provider {
	case ControllerK8sProvider:
		k8sClient, err := NewK8sClient(&conf.K8sProvider.KubeConfig)
		if err != nil {
			return nil, err
		}

		drainerConf := components.KubeNodeDrainerConf{
			NamespaceLabelSelector: conf.K8sProvider.NamespaceSelector,
			PodGracePeriod:         time.Duration(conf.GracePeriodSeconds) * time.Second,
			DryRun:                 conf.DryRun,
		}

		return components.NewKubeNodeDrainer(k8sClient, drainerConf, logger), nil
	}
	return nil, fmt.Errorf("node provider for type %s not found", conf.Provider)
}
