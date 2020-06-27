package ui

import (
	"log"
	"sync"

	"github.com/broothie/smock/ui/roundtrip"
	"github.com/gorilla/websocket"
)

type UI struct {
	Logger *log.Logger
	Port   int

	RoundTrips     map[string]*roundtrip.RoundTrip
	RoundTripsLock *sync.RWMutex

	Upgrader  websocket.Upgrader
	Conns     map[string]*websocket.Conn
	IsServing bool
}

func New(logger *log.Logger, port int) *UI {
	return &UI{
		Logger:         logger,
		Port:           port,
		RoundTrips:     make(map[string]*roundtrip.RoundTrip),
		RoundTripsLock: new(sync.RWMutex),
		Conns:          make(map[string]*websocket.Conn),
	}
}
