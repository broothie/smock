package ui

import (
	"log"
	"net/http"
	"time"

	"github.com/broothie/smock/pkg/id"
	"github.com/broothie/smock/pkg/reqlogger"
)

type UI struct {
	Server  *http.Server
	Port    int
	Logger  *log.Logger
	Entries map[string]Entry
}

func New(port int, logger *log.Logger) UI {
	ui := UI{
		Port:    port,
		Logger:  logger,
		Entries: make(map[string]Entry),
	}

	ui.setupServer()
	return ui
}

func (ui *UI) Start() {
	ui.Logger.Printf("ui @ http://localhost:%d\n", ui.Port)
	go ui.Server.ListenAndServe()
}

func (ui *UI) recordEntry(req *http.Request, res *http.Response, start, end time.Time) {
	id := id.FromContext(req.Context())

	entry, err := NewEntry(id, req, res, start, end)
	if err != nil {
		ui.Logger.Printf("failed to record entry for id=%s: %v", id, err)
		return
	}

	ui.Entries[id] = entry

	std := reqlogger.FormatFromReqRes(req, res, end.Sub(start))
	ui.Logger.Printf("%s | http://localhost:%d#%s", std, ui.Port, id)
}
