package ui

import (
	"net/http"
	"time"
)

func (ui *UI) Doer(req *http.Request) (*http.Response, error) {
	start := time.Now()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	end := time.Now()

	go ui.recordEntry(req, res, start, end)
	return res, err
}
