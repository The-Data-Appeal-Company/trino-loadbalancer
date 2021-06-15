package trino

import (
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNoopRetriever(t *testing.T) {
	retriever := Noop()
	stats, err := retriever.ClusterStatistics(models2.Coordinator{})
	require.NoError(t, err)
	require.Equal(t, stats, ClusterStatistics{})
}
