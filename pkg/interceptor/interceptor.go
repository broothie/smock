package interceptor

import (
	"net/http"
	"net/http/httptest"
)

type Interceptor struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func New(w http.ResponseWriter) *Interceptor {
	return &Interceptor{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (i *Interceptor) Write(body []byte) (int, error) {
	i.body = body
	return i.ResponseWriter.Write(body)
}

func (i *Interceptor) WriteHeader(statusCode int) {
	i.statusCode = statusCode
	i.ResponseWriter.WriteHeader(statusCode)
}

func (i *Interceptor) ToRecorder() *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	for key, values := range i.ResponseWriter.Header() {
		for _, value := range values {
			recorder.Header().Set(key, value)
		}
	}

	recorder.WriteHeader(i.statusCode)
	recorder.Write(i.body)
	return recorder
}

func (i *Interceptor) ToResponse() *http.Response {
	recorder := i.ToRecorder()
	res := recorder.Result()
	res.ContentLength = int64(recorder.Body.Len())
	return res
}
