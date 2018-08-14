package log

import "net/http"

func Middleware(width int, next http.HandlerFunc) http.HandlerFunc {
	return (&Logger{width}).Middleware(next)
}
