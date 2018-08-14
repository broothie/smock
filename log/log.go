package log

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

type Logger struct {
	Width int
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()
	dump, _ := httputil.DumpRequest(r, true)
	go func(dump string, t time.Time) {
		lines := Entrify([][]string{strings.Split(dump, "\r\n")}, l.Width)
		lines = append([]string{fmt.Sprintf("Request at %v", t)}, lines...)
		lines = append(lines, "\n")
		fmt.Println(strings.Join(lines, "\n"))
	}(string(dump), requestTime)
}

func (l *Logger) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.ServeHTTP(w, r)
		next(w, r)
	}
}

func (l *Logger) Wrap(handler http.Handler) http.Handler {
	l.Middleware(handler.ServeHTTP)
	return l
}

func Entrify(entries [][]string, width int) []string {
	var out []string

	divider := BuildDivider(width)
	out = append(out, divider)

	for _, entry := range entries {
		for _, bordered := range OutsideBorderify(entry, width) {
			out = append(out, bordered)
		}
		out = append(out, divider)
	}

	return out
}

func OutsideBorderify(in []string, width int) []string {
	var out []string
	lBar := "| %-"
	rBar := "s |"

	for _, line := range in {
		if len(line) == 0 {
			out = append(out, fmt.Sprintf(lBar+strconv.Itoa(width)+rBar, ""))
		}

		for j := 0; j < len(line); j += width {
			k := j + width

			left := lBar
			if j != 0 {
				left = "  %-"
			}

			right := rBar
			if width+j < len(line) {
				right = "s"
			} else {
				k = len(line)
			}

			out = append(out, fmt.Sprintf(
				left+strconv.Itoa(width)+right,
				line[j:k],
			))
		}
	}
	return out
}

func BuildDivider(width int) string {
	var dividerBuilder strings.Builder
	dividerBuilder.WriteString("+-")
	for i := 0; i < width; i++ {
		dividerBuilder.WriteRune('-')
	}
	dividerBuilder.WriteString("-+")
	return dividerBuilder.String()
}
