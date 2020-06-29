package handlers

import (
	"log"
	"net/http"
)

func Mock(logger *log.Logger, statusCode int, headers http.Header, body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for key, values := range headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(statusCode)
		if _, err := w.Write(body); err != nil {
			logger.Println(err)
		}
	}
}

func File(statusCode int, headers http.Header, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for key, values := range headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(statusCode)
		http.ServeFile(w, r, filename)
	}
}
