package proxy

import "net/http"

func Middleware(target string, log bool, width int, next http.HandlerFunc) http.HandlerFunc {
	return MustNewProxy(target, log, width).Middleware(next)
}

func To(target string, next http.HandlerFunc) http.HandlerFunc {
	return Middleware(target, false, 0, next)
}

func WithLoggingTo(target string, width int, next http.HandlerFunc) http.HandlerFunc {
	return Middleware(target, true, width, next)
}
