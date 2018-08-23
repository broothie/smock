package echo

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andydennisonbooth/smock/log"
)

type Echo struct {
	Width int
}

func New() *Echo {
	return &Echo{Width: log.DefaultWidth}
}

func (e *Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()
	entries := []string{fmt.Sprintf("Request at %v", requestTime)}

	defer func() {
		go func(entries []string) {
			fmt.Println(log.Entrify(e.Width, entries...))
		}(entries)
	}()

	dump, err := log.CleanDump(r)
	if err != nil {
		message := "Failed to dump request"
		entries = append(entries, message, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, message, err.Error())
		return
	}
	entries = append(entries, dump)

	fmt.Fprint(w, dump)
}
