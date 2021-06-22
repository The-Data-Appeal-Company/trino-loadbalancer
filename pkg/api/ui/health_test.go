package ui

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	rr := httptest.NewRecorder()
	healthProbe(rr, httptest.NewRequest(http.MethodGet, "http://localhost/healthz", nil))
	require.Equal(t, rr.Code, http.StatusOK)
}
