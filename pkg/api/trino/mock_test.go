package trino

import (
	"errors"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMockRetriever(t *testing.T) {

	mockErr := errors.New("generic err")
	mockStats := ClusterStatistics{}
	retriever := Mock(mockStats, mockErr)

	stats, err := retriever.ClusterStatistics(models2.Coordinator{})

	require.Equal(t, mockErr, err)
	require.Equal(t, mockStats, stats)

}
