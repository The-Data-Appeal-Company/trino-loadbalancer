package healthcheck

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestNoOpReturn(t *testing.T) {
	check := NoOp()

	testUrl, err := url.Parse("http://localhost:9876")
	require.NoError(t, err)

	result, err := check.Check(testUrl)
	require.NoError(t, err)
	require.Equal(t, true, result.IsAvailable())
}
