package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type proxyKey struct{}

var proxyk proxyKey

func ContextWithProxyKey(parent context.Context, reqRes *RequestResponse) context.Context {
	return context.WithValue(parent, proxyk, reqRes)
}

func GetReqResFromContext(ctx context.Context) *RequestResponse {
	reqRes, ok := ctx.Value(proxyk).(*RequestResponse)
	if !ok {
		return nil
	}

	return reqRes
}

type RequestResponse struct {
	Request  *http.Request
	Response *http.Response
}

func Proxy(target *url.URL) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetReq, err := copyRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		targetReq.Host = target.Host
		targetReq.URL, _ = url.Parse(target.String())
		targetReq.URL.Path = path.Join(target.Path, r.URL.Path)
		targetRes, err := http.DefaultClient.Do(targetReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := copyResponse(targetRes, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go func() {
			reqRes := GetReqResFromContext(r.Context())
			reqRes.Request = targetReq
			reqRes.Response = targetRes
		}()
	})
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
	for key, values := range res.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if _, err := w.Write(body); err != nil {
		return err
	}

	return nil
}
