package mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andydennisonbooth/smock/log"
)

type Mock struct {
	Reponse  string
	Width    int
	Filename string
}

func NewFromString(response string) *Mock {
	return &Mock{Reponse: response, Width: log.DefaultWidth}
}

func NewFromFile(filename string) (*Mock, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	m := NewFromString(string(bytes))
	m.Filename = filename
	return m, nil
}

func (m *Mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()
	entries := []string{fmt.Sprintf("Request at %v", requestTime)}

	defer func() {
		go func(entries []string) {
			fmt.Println(log.Entrify(m.Width, entries...))
		}(entries)
	}()

	dump, err := log.CleanDump(r)
	if err != nil {
		message := "Unable to dump request"
		entries = append(entries, message, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, message, err.Error())
		return
	}
	entries = append(entries, dump)

	if _, err := fmt.Fprint(w, m.Reponse); err == nil {
		if m.Filename == "" {
			entries = append(entries, fmt.Sprintf(
				"Responded with '%s'",
				m.Reponse,
			))
		} else {
			entries = append(entries, fmt.Sprintf(
				"Responded with contents of '%s'",
				m.Filename,
			))
		}
	} else {
		message := "Unable to write response"
		entries = append(entries, message)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, message, err.Error())
	}
}
