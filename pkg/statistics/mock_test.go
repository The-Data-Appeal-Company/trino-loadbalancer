package statistics

import (
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMockRetriever(t *testing.T) {

	mockErr := errors.New("generic err")
	mockStats := models.ClusterStatistics{}
	retriever := Mock(mockStats, mockErr)

	stats, err := retriever.ClusterStatistics(models.Coordinator{})

	require.Equal(t, mockErr, err)
	require.Equal(t, mockStats, stats)

}
