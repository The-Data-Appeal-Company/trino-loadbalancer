package configuration

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/notifier"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/controller/components"
	"github.com/go-redis/redis/v8"
	"time"
)

type ControllerConf struct {
	Features struct {
		SlowWorkerDrainer SlowWorkerDrainerConf `json:"slow_worker_drainer" yaml:"slow_worker_drainer" mapstructure:"slow_worker_drainer"`
	} `json:"features" yaml:"features" mapstructure:"features"`
}

type AutoscalerConf struct {
	Enabled    bool `yaml:"enabled" json:"enabled"`
	Kubernetes []struct {
		CoordinatorUri string        `json:"coordinatorUri,omitempty" yaml:"coordinatorUri,omitempty" json:"coordinatorUri,omitempty"`
		Namespace      string        `yaml:"namespace" json:"namespace,omitempty" json:"namespace,omitempty"`
		Deployment     string        `json:"deployment,omitempty" yaml:"deployment,omitempty" json:"deployment,omitempty"`
		Min            int           `json:"min,omitempty" yaml:"min,omitempty" json:"min,omitempty"`
		Max            int           `json:"max,omitempty" yaml:"max,omitempty" json:"max,omitempty"`
		ScaleAfter     time.Duration `json:"scaleAfter" yaml:"scaleAfter"`
	} `yaml:"kubernetes" json:"kubernetes"`
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
	Analyzer struct {
		StdDeviationRation float64 `json:"std_deviation_ratio" yaml:"std_deviation_ratio" mapstructure:"std_deviation_ratio"`
	} `json:"analyzer" yaml:"analyzer" mapstructure:"analyzer"`
}

func CreateHandlers(redisClient redis.UniversalClient, logger logging.Logger, notifier notifier.Notifier, conf ControllerConf) (components.QueryHandler, error) {
	handlers := make([]components.QueryHandler, 0)

	slowNodeDrainerConf := conf.Features.SlowWorkerDrainer
	if slowNodeDrainerConf.Enabled {
		slowWorkerHandler, err := createDrainSlowWorkerNodeHandler(redisClient, logger, notifier, slowNodeDrainerConf)
		if err != nil {
			return nil, err
		}

		handlers = append(handlers, slowWorkerHandler)
	}

	return components.NewMultiQueryHandler(handlers...), nil
}

func createDrainSlowWorkerNodeHandler(redisClient redis.UniversalClient, logger logging.Logger, notifier notifier.Notifier, slowNodeDrainerConf SlowWorkerDrainerConf) (*components.SlowNodeDrainer, error) {
	analyzer := components.NewTrinoSlowNodeAnalyzer(components.TrinoSlowNodeAnalyzerConfig{
		StdDeviationRatio: slowNodeDrainerConf.Analyzer.StdDeviationRation,
	})

	nodeDrainer, err := createNodeDrainer(slowNodeDrainerConf, logger)
	if err != nil {
		return nil, err
	}
	slowNodeMarker := components.NewRedisSlowNodeMarker(redisClient)
	return components.NewSlowNodeDrainer(analyzer, nodeDrainer, slowNodeMarker, components.SlowNodeDrainerConf{DrainThreshold: slowNodeDrainerConf.DrainThreshold}, logger, notifier), nil
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
