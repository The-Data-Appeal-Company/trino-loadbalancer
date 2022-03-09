package trino

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestNoopRetriever(t *testing.T) {
	retriever := Noop()
	stats, err := retriever.ClusterStatistics(&url.URL{})
	require.NoError(t, err)
	require.Equal(t, stats, ClusterStatistics{})
}
