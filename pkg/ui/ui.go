package ui

import (
	"log"
	"net/http"
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
