package id

import (
	"context"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type key struct{}

var k key

func FromContext(ctx context.Context) string {
	return ctx.Value(k).(string)
}

func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, k, id)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(WithContext(r.Context(), randID(6))))
	})
}

func randID(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	runes := []rune(chars)

	id := make([]rune, length)
	for i := range id {
		id[i] = runes[rand.Intn(len(runes))]
	}

	return string(id)
}
