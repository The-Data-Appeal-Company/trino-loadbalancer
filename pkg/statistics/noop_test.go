package statistics

import (
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNoopRetriever(t *testing.T) {
	retriever := Noop()
	stats, err := retriever.GetStatistics(models.Coordinator{})
	require.NoError(t, err)
	require.Equal(t, stats, models.ClusterStatistics{})
}
