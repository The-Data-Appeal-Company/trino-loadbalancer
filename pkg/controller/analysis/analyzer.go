package analysis

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/montanaflynn/stats"
	"time"
)

const StdDeviationRatio = 1.1

type NodeAnalyzer interface {
	Analyze(trino.QueryDetail) ([]string, error)
}

type TrinoNodeAnalyzer struct {
}

func (a TrinoNodeAnalyzer) Analyze(queryDetail trino.QueryDetail) ([]string, error) {
	tasks := make([]trino.Tasks, 0)
	for _, stage := range queryDetail.OutputStage.SubStages {
		tasks = append(tasks, stage.Tasks...)
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
		taskTimes = append(taskTimes, float64(duration.Microseconds()))
	}

	averageTime, err := stats.Mean(taskTimes)
	if err != nil {
		return nil, err
	}
	stdDeviationTime, err := stats.StandardDeviation(taskTimes)
	if err != nil {
		return nil, err
	}

	nodesToEvict := make([]string, 0)
	nodeEvictionThreshold := averageTime + stdDeviationTime*StdDeviationRatio
	for nodeId, duration := range taskMap {
		if float64(duration.Microseconds()) > nodeEvictionThreshold {
			nodesToEvict = append(nodesToEvict, nodeId)
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
