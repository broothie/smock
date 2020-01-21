package ui

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/broothie/smock/pkg/interceptor"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (ui *UI) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		interceptor := interceptor.New(w)
		start := time.Now()
		next.ServeHTTP(interceptor, r)
		end := time.Now()

		go ui.recordEntry(r, interceptor.ToResponse(), start, end)
	})
}
