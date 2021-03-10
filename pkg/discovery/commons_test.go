package discovery

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContains(t *testing.T) {

	require.True(t, containsAll(map[string]string{}, map[string]string{}))

	require.True(t, containsAll(map[string]string{
		"a": "a",
	}, map[string]string{
		"a": "a",
	}))

	require.True(t, containsAll(map[string]string{
		"a": "a",
		"b": "a",
	}, map[string]string{
		"a": "a",
	}))

	require.False(t, containsAll(map[string]string{
		"a": "a",
	}, map[string]string{
		"a": "a",
		"b": "a",
	}))

	require.True(t, containsAll(map[string]string{
		"a": "a",
	}, map[string]string{}))
}
