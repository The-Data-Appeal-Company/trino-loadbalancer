package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Interceptor interface {
	Handle(*http.Request, *http.Response) error
}

type HttpProxy interface {
	Handle(http.ResponseWriter, *http.Request) error
}

type ReverseProxy struct {
	proxy *httputil.ReverseProxy
}

func NewReverseProxy(target *url.URL, interceptor Interceptor) *ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ModifyResponse = func(response *http.Response) error {
		if interceptor == nil {
			return nil
		}
		return interceptor.Handle(response.Request, response)
	}

	return &ReverseProxy{
		proxy: proxy,
	}
}

func (r ReverseProxy) Handle(w http.ResponseWriter, req *http.Request) error {
	r.proxy.ServeHTTP(w, req)
	return nil
}
