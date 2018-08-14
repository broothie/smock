package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/andydennisonbooth/smock/log"
)

type proxy struct {
	targetURL *url.URL
	client    *http.Client
	log       bool
	width     int
}

func MustNewProxy(target string, log bool, width int) *proxy {
	p, err := NewProxy(target, log, width)
	if err != nil {
		panic(err)
	}
	return p
}

func NewProxy(target string, log bool, width int) (*proxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	return &proxy{targetURL, http.DefaultClient, log, width}, nil
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()

	actualRequest, err := copyRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Unable to read request")
		return
	}
	actualRequest.Host = p.targetURL.Host
	actualRequest.URL, _ = url.Parse(p.targetURL.String())
	actualRequest.URL.Path = path.Join(p.targetURL.Path, r.URL.Path)
	reqDump, _ := httputil.DumpRequestOut(actualRequest, true)

	response, err := p.client.Do(actualRequest)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Request to external server failed")
		return
	}

	if p.log {
		resDump, _ := httputil.DumpResponse(response, true)
		go func(reqDump, resDump string, t time.Time) {
			reqRes := log.Entrify([][]string{
				strings.Split(reqDump, "\r\n"),
				strings.Split(resDump, "\r\n"),
			}, p.width)
			reqRes = append([]string{fmt.Sprintf("Request/response at %v", t)}, reqRes...)
			reqRes = append(reqRes, "\n")
			fmt.Println(strings.Join(reqRes, "\n"))
		}(string(reqDump), string(resDump), requestTime)
	}

	if err := copyResponse(response, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Writing external server response failed")
		return
	}
}

func (p *proxy) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
		next(w, r)
	}
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
