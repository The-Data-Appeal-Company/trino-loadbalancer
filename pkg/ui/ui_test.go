package ui

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeStaticFile(t *testing.T) {
	ui := New("testdata")

	r, err := http.NewRequest("GET", "/ui/script.js", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	ui.Router().ServeHTTP(w, r)

	loc := w.Result().Header.Get("Location")

	require.Equal(t, http.StatusOK, w.Result().StatusCode)

	responseBody, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)

	expectedBody, err := ioutil.ReadFile("testdata/script.js")
	require.NoError(t, err)

	require.Equal(t, string(expectedBody), string(responseBody))

}
