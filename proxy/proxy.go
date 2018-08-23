package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"time"

	"github.com/andydennisonbooth/smock/log"
)

type Proxy struct {
	TargetURL *url.URL
	Client    *http.Client
	Width     int
}

func New(target string, client *http.Client) (*Proxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = http.DefaultClient
	}

	return &Proxy{
		TargetURL: targetURL,
		Client:    client,
		Width:     log.DefaultWidth,
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()
	entries := []string{fmt.Sprintf("Request at %v", requestTime)}

	defer func() {
		go func(entries []string) {
			fmt.Println(log.Entrify(p.Width, entries...))
		}(entries)
	}()

	requestDump, err := log.CleanDump(r)
	if err != nil {
		message := "Unable to dump original request"
		entries = append(entries, message, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, message, err.Error())
		return
	}
	entries = append(entries, requestDump)

	// Copy request
	requestToMake, err := copyRequest(r)
	if err != nil {
		message := "Unable to copy original request"
		entries = append(entries, message, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, message, err.Error())
		return
	}

	// Reformat request to use proxy settings
	requestToMake.Host = p.TargetURL.Host
	requestToMake.URL, _ = url.Parse(p.TargetURL.String())
	requestToMake.URL.Path = path.Join(p.TargetURL.Path, r.URL.Path)

	// Make request to proxied server
	response, err := p.Client.Do(requestToMake)
	responseTime := time.Now()
	if err != nil {
		message := "Request to external server failed"
		entries = append(entries, message, err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, message, err.Error())
		return
	}
	entries = append(entries, fmt.Sprintf("Response at %v", responseTime))

	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		message := "Unable to dump external server response"
		entries = append(entries, message, err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, message, err.Error())
		return
	}
	entries = append(entries, string(responseDump))

	if err := copyResponse(response, w); err != nil {
		message := "Writing external server response failed"
		entries = append(entries, message, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, message, err.Error())
		return
	}
}

func (p *Proxy) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
		next(w, r)
	}
}

func (p *Proxy) Entrify(entries ...string) string {
	return log.Entrify(p.Width, entries...)
}

func copyRequest(r *http.Request) (*http.Request, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return &http.Request{
		Method:           r.Method,
		URL:              r.URL,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           r.Header,
		Body:             ioutil.NopCloser(bytes.NewBuffer(body)),
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Close:            r.Close,
		Host:             r.Host,
		Form:             r.Form,
		PostForm:         r.PostForm,
		MultipartForm:    r.MultipartForm,
	}, nil
}

func copyResponse(res *http.Response, w http.ResponseWriter) error {
	for name, header := range res.Header {
		for _, value := range header {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	w.Write(body)

	return nil
}
