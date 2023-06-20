package serving

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServeStaticFile(t *testing.T) {
	ui := New("testdata")

	r, err := http.NewRequest("GET", "/ui/script.js", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	ui.Router().ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	responseBody, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)

	expectedBody, err := os.ReadFile("testdata/script.js")
	require.NoError(t, err)

	require.Equal(t, string(expectedBody), string(responseBody))

}
