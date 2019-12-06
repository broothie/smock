package ui

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/broothie/smock/pkg/handlers"
	"github.com/broothie/smock/pkg/interceptor"
	"github.com/broothie/smock/pkg/reqlogger"
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
		id := randID(6)
		start := time.Now()
		next.ServeHTTP(interceptor, r.WithContext(handlers.ContextWithProxyKey(r.Context(), new(handlers.RequestResponse))))
		end := time.Now()

		// Record entry and log
		go func() {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			entry, err := NewEntry(id, start, end, r, interceptor)
			if err != nil {
				ui.Logger.Printf("failed to recorder entry for id=%s", id)
				return
			}

			reqRes := handlers.GetReqResFromContext(r.Context())
			if reqRes != nil {
				entry.TargetRequest = Request{}
				entry.TargetResponse = Response{}
			}

			ui.Entries[id] = entry
			std := reqlogger.FormatFromReqRes(r, interceptor.ToRecorder(), end.Sub(start))
			ui.Logger.Printf("%s | http://localhost:%d#%s", std, ui.Port, id)
		}()
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
