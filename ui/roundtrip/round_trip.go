package roundtrip

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type RoundTrip struct {
	ID          string
	StartTime   time.Time
	EndTime     time.Time
	RawRequest  []byte
	RawResponse []byte

	Request      *http.Request
	RequestBody  []byte
	Response     *http.Response
	ResponseBody []byte
}

func New(id string, startTime, endTime time.Time, rawRequest, rawResponse []byte) (*RoundTrip, error) {
	roundTrip := &RoundTrip{
		ID:          id,
		StartTime:   startTime,
		EndTime:     endTime,
		RawRequest:  rawRequest,
		RawResponse: rawResponse,
	}

	if err := roundTrip.setup(); err != nil {
		return nil, err
	}

	return roundTrip, nil
}

func (rt *RoundTrip) setup() error {
	var err error
	rt.Request, err = http.ReadRequest(bufio.NewReader(bytes.NewBuffer(rt.RawRequest)))
	if err != nil {
		return err
	}

	rt.Response, err = http.ReadResponse(bufio.NewReader(bytes.NewBuffer(rt.RawResponse)), rt.Request)
	if err != nil {
		return err
	}

	rt.RequestBody, err = ioutil.ReadAll(rt.Request.Body)
	if err != nil {
		return err
	}

	rt.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rt.RequestBody[:]))

	rt.ResponseBody, err = ioutil.ReadAll(rt.Response.Body)
	if err != nil {
		return err
	}

	rt.Response.Body = ioutil.NopCloser(bytes.NewBuffer(rt.ResponseBody[:]))
	return nil
}

func (rt *RoundTrip) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id": rt.ID,
		"request": map[string]interface{}{
			"method":        rt.Request.Method,
			"path":          rt.Request.URL.Path,
			"query":         rt.Request.URL.RawQuery,
			"headers":       rt.Request.Header,
			"contentLength": rt.Request.ContentLength,
			"body":          string(rt.RequestBody),
			"raw":           string(rt.RawRequest),
		},
		"response": map[string]interface{}{
			"code":          rt.Response.StatusCode,
			"status":        http.StatusText(rt.Response.StatusCode),
			"headers":       rt.Response.Header,
			"contentLength": rt.Response.ContentLength,
			"body":          string(rt.ResponseBody),
			"raw":           string(rt.RawResponse),
		},
	})
}
