package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/broothie/smock/ui/assets"
	"github.com/broothie/smock/ui/roundtrip"
	"github.com/gorilla/websocket"
)

func (ui *UI) Serve() {
	ui.IsServing = true
	ui.Logger.Printf("serving ui server @ %d\n", ui.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", ui.Port), ui.Handler()); err != nil {
		ui.Logger.Println(err)
		os.Exit(1)
	}
}

func (ui *UI) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/round_trips", ui.getRoundTrips)
	mux.HandleFunc("/round_trip", ui.getRoundTrip)
	mux.HandleFunc("/ws", ui.ws)

	return mux
}

func index(w http.ResponseWriter, _ *http.Request) {
	if _, err := fmt.Fprint(w, assets.IndexHTML); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ui *UI) getRoundTrips(w http.ResponseWriter, _ *http.Request) {
	trips := make([]*roundtrip.RoundTrip, len(ui.RoundTrips))
	counter := 0
	for _, roundTrip := range ui.RoundTrips {
		trips[counter] = roundTrip
		counter++
	}

	sort.Slice(trips, func(i, j int) bool { return trips[i].StartTime.Before(trips[j].StartTime) })

	body, err := json.Marshal(trips)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ui *UI) getRoundTrip(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	trip, found := ui.FetchRoundTrip(id)
	if !found {
		http.Error(w, fmt.Sprintf("no roundtrip found for <%s>", id), http.StatusNotFound)
		return
	}

	body, err := json.Marshal(trip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ui *UI) ws(w http.ResponseWriter, r *http.Request) {
	conn, err := ui.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := randID()
	ui.Conns[id] = conn
	defer func() {
		delete(ui.Conns, id)
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if _, isClose := err.(*websocket.CloseError); isClose {
				break
			}
		}
	}
}

func (ui *UI) alertConnections() {
	for _, connection := range ui.Conns {
		if err := connection.WriteMessage(websocket.TextMessage, []byte{}); err != nil {
			fmt.Println("failed to ping")
		}
	}
}
