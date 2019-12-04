package reqlogger

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/broothie/smock/pkg/interceptor"
)

func New(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			interceptor := interceptor.New(w)

			before := time.Now()
			next.ServeHTTP(interceptor, r)
			since := time.Since(before)

			logger.Println(FormatFromReqRes(r, interceptor.ToRecorder(), since))
		})
	}
}

func Format(method, path, rawQuery string, reqContentLength, code, resContentLength int, since time.Duration) string {
	return fmt.Sprintf("%s %s%s %dB | %d %s %dB | %v",
		method,
		path,
		rawQuery,
		reqContentLength,
		code,
		http.StatusText(code),
		resContentLength,
		since,
	)
}

func FormatFromReqRes(req *http.Request, recorder *httptest.ResponseRecorder, since time.Duration) string {
	return Format(
		req.Method,
		req.URL.Path,
		req.URL.RawQuery,
		int(req.ContentLength),
		recorder.Code,
		recorder.Body.Len(),
		since,
	)
}
