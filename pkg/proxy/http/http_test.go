package http

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttpProxy(t *testing.T) {
	backendSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, err := w.Write([]byte("ok"))
		require.NoError(t, err)
	}))
	defer backendSrv.Close()

	uri, err := url.Parse(backendSrv.URL)
	require.NoError(t, err)

	proxy := NewReverseProxy(uri, nil)

	proxySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.proxy.ServeHTTP(w, r)
	}))
	defer proxySrv.Close()

	resp, err := http.Get(proxySrv.URL)
	require.NoError(t, err)

	require.Equal(t, http.StatusAccepted, resp.StatusCode)

}

type TestableInterceptor struct {
	calls int
}

func newTestableInterceptor() *TestableInterceptor {
	return &TestableInterceptor{}
}

func (t *TestableInterceptor) Handle(request *http.Request, response *http.Response) error {
	t.calls++
	return nil
}

func TestHttpProxyWithInterceptor(t *testing.T) {
	interceptor := newTestableInterceptor()

	backendSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, err := w.Write([]byte("ok"))
		require.NoError(t, err)
	}))
	defer backendSrv.Close()

	uri, err := url.Parse(backendSrv.URL)
	require.NoError(t, err)

	proxy := NewReverseProxy(uri, interceptor)

	proxySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.proxy.ServeHTTP(w, r)
	}))
	defer proxySrv.Close()

	resp, err := http.Get(proxySrv.URL)
	require.NoError(t, err)

	require.Equal(t, http.StatusAccepted, resp.StatusCode)
	require.Equal(t, 1, interceptor.calls)
}
