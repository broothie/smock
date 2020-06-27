package ui

import (
	"math/rand"
	"time"

	"github.com/broothie/smock/ui/roundtrip"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (ui *UI) AddRoundTrip(id string, startTime, endTime time.Time, rawRequest, rawResponse []byte) {
	roundTrip, err := roundtrip.New(id, startTime, endTime, rawRequest, rawResponse)
	if err != nil {
		ui.Logger.Printf("failed to add roundtrip <%s>: %v\n", id, err)
		return
	}

	defer func() {
		ui.Log(id)
		if ui.IsServing {
			ui.alertConnections()
		}
	}()

	ui.RoundTripsLock.Lock()
	defer ui.RoundTripsLock.Unlock()
	ui.RoundTrips[id] = roundTrip
}

func (ui *UI) FetchRoundTrip(id string) (*roundtrip.RoundTrip, bool) {
	ui.RoundTripsLock.RLock()
	roundTrip, found := ui.RoundTrips[id]
	ui.RoundTripsLock.RUnlock()

	return roundTrip, found
}

func (ui *UI) newID() string {
	id := randID()
	for _, idExists := ui.RoundTrips[id]; idExists; _, idExists = ui.RoundTrips[id] {
		id = randID()
	}

	return id
}

func randID() string {
	return randString(6)
}

func randString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	runes := []rune(chars)

	id := make([]rune, length)
	for i := range id {
		id[i] = runes[rand.Intn(len(runes))]
	}

	return string(id)
}
