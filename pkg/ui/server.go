package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

func (ui *UI) setupServer() {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, html); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	serveMux.HandleFunc("/entries", func(w http.ResponseWriter, r *http.Request) {
		if id := r.URL.Query().Get("id"); id != "" {
			bytes, err := json.Marshal(ui.Entries[id])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := w.Write(bytes); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		var entries []Entry
		for _, entry := range ui.Entries {
			entries = append(entries, entry)
		}

		if len(entries) == 0 {
			entries = []Entry{}
		}

		sort.Slice(entries, func(i, j int) bool { return entries[i].Start.Before(entries[j].Start) })

		bytes, err := json.Marshal(entries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(bytes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	ui.Server = &http.Server{Addr: fmt.Sprintf(":%d", ui.Port), Handler: serveMux}
}
