package healthcheck

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNotAvailableOnDefault(t *testing.T) {
	h := Health{}
	require.Equal(t, false, h.IsAvailable())
}

func TestAvailableWhenHealthy(t *testing.T) {
	h := Health{Status: StatusHealthy}
	require.Equal(t, true, h.IsAvailable())
}

func TestNotAvailableWhenUnknown(t *testing.T) {
	h := Health{Status: StatusUnknown}
	require.Equal(t, false, h.IsAvailable())
}
