package ui

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/broothie/smock/pkg/interceptor"
	"github.com/broothie/smock/pkg/reqlogger"
)

type Entry struct {
	ID       string        `json:"id"`
	Line     string        `json:"line"`
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Elapsed  time.Duration `json:"elapsed"`
	Request  Request       `json:"request"`
	Response Response      `json:"response"`
}

type Request struct {
	Method   string `json:"method"`
	Path     string `json:"path"`
	Query    string `json:"query"`
	Protocol string `json:"protocol"`
	Raw      string `json:"raw"`
}

type Response struct {
	Code     int    `json:"code"`
	Protocol string `json:"protocol"`
	Raw      string `json:"raw"`
}

func NewEntry(id string, start, end time.Time, req *http.Request, interceptor *interceptor.Interceptor) (Entry, error) {
	recorder := interceptor.ToRecorder()
	res := recorder.Result()

	rawReq, err := httputil.DumpRequest(req, true)
	if err != nil {
		return Entry{}, err
	}

	rawRes, err := httputil.DumpResponse(res, true)
	if err != nil {
		return Entry{}, err
	}

	line := fmt.Sprintf("%s | %s", id, reqlogger.FormatFromReqRes(req, recorder, end.Sub(start)))
	return Entry{
		ID:      id,
		Line:    line,
		Start:   start,
		End:     end,
		Elapsed: end.Sub(start),
		Request: Request{
			Method:   req.Method,
			Path:     req.URL.Path,
			Query:    req.URL.RawQuery,
			Protocol: req.Proto,
			Raw:      string(rawReq),
		},
		Response: Response{
			Code:     recorder.Code,
			Protocol: res.Proto,
			Raw:      string(rawRes),
		},
	}, nil
}
