package reqlogger

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/broothie/smock/pkg/interceptor"
)

func Wrap(next http.Handler, logger *log.Logger) http.Handler {
	return New(logger)(next)
}

func New(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			interceptor := interceptor.New(w)

			before := time.Now()
			next.ServeHTTP(interceptor, r)
			since := time.Since(before)

			logger.Println(FormatFromReqRes(r, interceptor.ToResponse(), since))
		})
	}
}

func Format(method, path, rawQuery string, reqContentLength, code, resContentLength int, elapsed time.Duration) string {
	questionMark := "?"
	if rawQuery == "" {
		questionMark = ""
	}

	return fmt.Sprintf("%s %s%s%s %dB | %d %s %dB | %v",
		method,
		path,
		questionMark,
		rawQuery,
		reqContentLength,
		code,
		http.StatusText(code),
		resContentLength,
		elapsed,
	)
}

func FormatFromReqRes(req *http.Request, res *http.Response, elapsed time.Duration) string {
	return Format(
		req.Method,
		req.URL.Path,
		req.URL.RawQuery,
		int(req.ContentLength),
		res.StatusCode,
		int(res.ContentLength),
		elapsed,
	)
}
