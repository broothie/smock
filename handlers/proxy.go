package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Proxy(url *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reverseProxy := httputil.NewSingleHostReverseProxy(url)
		reverseProxy.Transport = transport{host: url.Host}
		reverseProxy.ServeHTTP(w, r)
	}
}

type transport struct {
	host string
}

func (t transport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Host = t.host
	r.Header.Del("X-Forwarded-For")
	return http.DefaultTransport.RoundTrip(r)
}
