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
	analyzer := NewTrinoSlowNodeAnalyzer(TrinoSlowNodeAnalyzerConfig{StdDeviationRatio: 1.1})
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
								ElapsedTime: "22ms",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-02",
							},
							Stats: trino.Stats{
								ElapsedTime: "2ms",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-03",
							},
							Stats: trino.Stats{
								ElapsedTime: "20ms",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-04",
							},
							Stats: trino.Stats{
								ElapsedTime: "1ms",
							},
						},
						{
							TaskStatus: trino.TaskStatus{
								NodeID: "node-05",
							},
							Stats: trino.Stats{
								ElapsedTime: "1ms",
							},
						},
					},
				},
			},
		},
	}
	analyzer := NewTrinoSlowNodeAnalyzer(TrinoSlowNodeAnalyzerConfig{StdDeviationRatio: 1.1})
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
	analyzer := NewTrinoSlowNodeAnalyzer(TrinoSlowNodeAnalyzerConfig{StdDeviationRatio: 1.1})
	nodes, err := analyzer.Analyze(queryDetail)
	require.NoError(t, err)
	require.Len(t, nodes, 0)
}

func TestShouldDetectNoNodesIfNoTaskAreAvailable(t *testing.T) {
	queryDetail := trino.QueryDetail{
		OutputStage: trino.OutputStage{
			SubStages: nil,
		},
	}
	analyzer := NewTrinoSlowNodeAnalyzer(TrinoSlowNodeAnalyzerConfig{StdDeviationRatio: 1.1})
	nodes, err := analyzer.Analyze(queryDetail)
	require.NoError(t, err)
	require.Len(t, nodes, 0)
}
