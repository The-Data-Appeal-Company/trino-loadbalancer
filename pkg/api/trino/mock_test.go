package trino

import (
	"errors"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestMockRetriever(t *testing.T) {

	mockErr := errors.New("generic err")
	mockStats := ClusterStatistics{}
	retriever := Mock(mockStats, mockErr)

	stats, err := retriever.ClusterStatistics(&url.URL{})

	require.Equal(t, mockErr, err)
	require.Equal(t, mockStats, stats)

}
