package components

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShouldDetectSingleSlowNode(t *testing.T) {
	queryDetail := trino.QueryDetail{
		OutputStage: trino.OutputStage{
			SubStages: []trino.SubStages{
				{
					Tasks: []trino.Tasks{
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-01",
							},
							Stats: trino.Stats{
								ElapsedTime: "2s",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-02",
							},
							Stats: trino.Stats{
								ElapsedTime: "2s",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-03",
							},
							Stats: trino.Stats{
								ElapsedTime: "2s",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-04",
							},
							Stats: trino.Stats{
								ElapsedTime: "100s",
							},
						},
					},
				},
			},
		},
	}
	analyzer := TrinoSlowNodeAnalyzer{}
	nodes, err := analyzer.Analyze(queryDetail)
	require.NoError(t, err)
	require.Len(t, nodes, 1)
	require.Equal(t, nodes[0], SlowNodeRef{NodeID: "node-04"})
}

func TestShouldDetectMultiplesSlowNode(t *testing.T) {
	queryDetail := trino.QueryDetail{
		OutputStage: trino.OutputStage{
			SubStages: []trino.SubStages{
				{
					Tasks: []trino.Tasks{
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-01",
							},
							Stats: trino.Stats{
								ElapsedTime: "22us",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-02",
							},
							Stats: trino.Stats{
								ElapsedTime: "2us",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-03",
							},
							Stats: trino.Stats{
								ElapsedTime: "20us",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-04",
							},
							Stats: trino.Stats{
								ElapsedTime: "1us",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-05",
							},
							Stats: trino.Stats{
								ElapsedTime: "1us",
							},
						},
					},
				},
			},
		},
	}
	analyzer := TrinoSlowNodeAnalyzer{}
	nodes, err := analyzer.Analyze(queryDetail)
	require.NoError(t, err)
	require.Len(t, nodes, 2)
	require.Contains(t, nodes, SlowNodeRef{NodeID: "node-01"})
	require.Contains(t, nodes, SlowNodeRef{NodeID: "node-03"})
}

func TestShouldDetectNoNodes(t *testing.T) {
	queryDetail := trino.QueryDetail{
		OutputStage: trino.OutputStage{
			SubStages: []trino.SubStages{
				{
					Tasks: []trino.Tasks{
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-01",
							},
							Stats: trino.Stats{
								ElapsedTime: "10ms",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-02",
							},
							Stats: trino.Stats{
								ElapsedTime: "8ms",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-03",
							},
							Stats: trino.Stats{
								ElapsedTime: "11ms",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-04",
							},
							Stats: trino.Stats{
								ElapsedTime: "13mss",
							},
						},
					},
				},
			},
		},
	}
	analyzer := TrinoSlowNodeAnalyzer{}
	nodes, err := analyzer.Analyze(queryDetail)
	require.NoError(t, err)
	require.Len(t, nodes, 0)
}
