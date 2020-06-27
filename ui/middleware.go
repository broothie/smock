package ui

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"time"
)

func (ui *UI) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := ui.newID()
		recorder := httptest.NewRecorder()

		// Dump request
		rawRequest, err := httputil.DumpRequest(r, r.ContentLength > 0)
		if err != nil {
			ui.Logger.Printf("failed to dump request for roundtrip <%s>: %v", id, err)
		}

		// Let handler do its thing
		startTime := time.Now()
		next.ServeHTTP(recorder, r)
		endTime := time.Now()

		// Dump response
		response := recorder.Result()
		response.ContentLength = int64(recorder.Body.Len())
		rawResponse, err := httputil.DumpResponse(response, response.ContentLength > 0)
		if err != nil {
			ui.Logger.Printf("failed to dump response for roundtrip <%s>: %v", id, err)
		}

		// Add to history
		if err == nil {
			go ui.AddRoundTrip(id, startTime, endTime, rawRequest, rawResponse)
		}

		// Write actual response
		for key, values := range recorder.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
	})
}
