package autoscaler

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/configuration"
	"math"
)

const (
	fixedScaleUpStrategy       string = "fixed"
	linearScaleUpStrategy      string = "linear"
	logarithmicScaleUpStrategy string = "logarithmic"
)

func calculateScaleUpWorker(currentWorker, activeUser int, strategy configuration.AutoscalerScaleUpStrategy) (int, error) {

	switch strategy.Type {
	case fixedScaleUpStrategy:
		return strategy.Max, nil
	case linearScaleUpStrategy:
		return maxInt(currentWorker, maxInt(currentWorker+(activeUser*strategy.WorkerPerUser), strategy.Max)), nil
	case logarithmicScaleUpStrategy:
		return maxInt(currentWorker, maxInt(currentWorker+int(math.Log(float64(activeUser)*float64(strategy.WorkerPerUser))), strategy.Max)), nil
	default:
		return 0, fmt.Errorf("invalid scale up strategy")
	}
}

func maxInt(value1, value2 int) int {
	if value1 > value2 {
		return value1
	} else {
		return value2
	}
}
