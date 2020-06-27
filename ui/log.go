package ui

import (
	"fmt"
	"net/http"
	"strings"
)

func (ui *UI) Log(id string) {
	roundTrip, found := ui.FetchRoundTrip(id)
	if !found {
		return
	}

	req := roundTrip.Request
	res := roundTrip.Response

	builder := new(strings.Builder)
	builder.WriteString(req.Method + " ")
	builder.WriteString(req.URL.Path)
	if req.URL.RawQuery != "" {
		builder.WriteString("?" + req.URL.RawQuery)
	}

	builder.WriteString(" ")

	if req.ContentLength > 0 {
		builder.WriteString(fmt.Sprintf("%dB ", req.ContentLength))
	}

	builder.WriteString(fmt.Sprintf("| %d %s ", res.StatusCode, http.StatusText(res.StatusCode)))

	if res.ContentLength > 0 {
		builder.WriteString(fmt.Sprintf("%dB ", res.ContentLength))
	}

	builder.WriteString(fmt.Sprintf("| %v", roundTrip.EndTime.Sub(roundTrip.StartTime)))

	if ui.IsServing {
		builder.WriteString(fmt.Sprintf(" | http://localhost:%d?id=%s", ui.Port, id))
	}

	ui.Logger.Println(builder.String())
}
