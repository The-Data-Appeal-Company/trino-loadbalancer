package ui

import (
	logging2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiHealth(t *testing.T) {

	api := NewApi(nil, nil, nil, logging2.Noop())

	r, err := http.NewRequest("GET", "/api/health", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	api.Router().ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
}
