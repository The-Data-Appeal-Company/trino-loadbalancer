package components

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/montanaflynn/stats"
	"time"
)

type SlowNodeRef struct {
	NodeID string
}

type SlowNodeAnalyzer interface {
	Analyze(trino.QueryDetail) ([]SlowNodeRef, error)
}

type TrinoSlowNodeAnalyzerConfig struct {
	StdDeviationRatio float64
}

type TrinoSlowNodeAnalyzer struct {
	conf TrinoSlowNodeAnalyzerConfig
}

func NewTrinoSlowNodeAnalyzer(conf TrinoSlowNodeAnalyzerConfig) TrinoSlowNodeAnalyzer {
	return TrinoSlowNodeAnalyzer{conf: conf}
}

func (a TrinoSlowNodeAnalyzer) Analyze(queryDetail trino.QueryDetail) ([]SlowNodeRef, error) {
	tasks := make([]trino.Tasks, 0)
	for _, stage := range queryDetail.OutputStage.SubStages {
		tasks = append(tasks, stage.Tasks...)
	}

	if len(tasks) == 0 {
		return nil, nil
	}

	taskMap := make(map[string]time.Duration)
	for _, task := range tasks {
		val, found := taskMap[task.TaskStatus.NodeID]
		if !found || compareTasksElapsedTime(val, task) {
			taskMap[task.TaskStatus.NodeID] = getDurationFromTask(task)
		}
	}

	taskTimes := make([]float64, 0)
	for _, duration := range taskMap {
		taskTimes = append(taskTimes, float64(duration.Milliseconds()))
	}

	averageTime, err := stats.Mean(taskTimes)
	if err != nil {
		return nil, err
	}
	stdDeviationTime, err := stats.StandardDeviation(taskTimes)
	if err != nil {
		return nil, err
	}

	nodesToEvict := make([]SlowNodeRef, 0)
	nodeEvictionThreshold := averageTime + stdDeviationTime*a.conf.StdDeviationRatio
	for nodeId, duration := range taskMap {
		if float64(duration.Milliseconds()) > nodeEvictionThreshold {
			nodesToEvict = append(nodesToEvict, SlowNodeRef{NodeID: nodeId})
		}
	}

	return nodesToEvict, nil
}

func compareTasksElapsedTime(first time.Duration, second trino.Tasks) bool {
	return first < getDurationFromTask(second)
}

func getDurationFromTask(task trino.Tasks) time.Duration {
	elapsedTime, err := task.GetElapsedTime()
	if err != nil {
		elapsedTime = time.Duration(0)
	}
	return elapsedTime
}
