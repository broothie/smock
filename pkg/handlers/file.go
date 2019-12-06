package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
)

func File(logger *log.Logger, filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			logger.Println(err)
			return
		}

		Mock(logger, http.StatusOK, http.Header{}, contents).ServeHTTP(w, r)
	}
}
