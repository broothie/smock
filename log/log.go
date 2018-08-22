package log

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var newlineRegexp = regexp.MustCompile(`\r\n`)

type Logger struct {
	Width int
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()
	dump, _ := httputil.DumpRequest(r, true)
	go func(dump string, t time.Time) {
		fmt.Println(l.Entrify(fmt.Sprintf("Request at %v", t), dump))
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

func (l *Logger) Entrify(entries ...string) string {
	return Entrify(entries, l.Width)
}

func Entrify(entries []string, width int) string {
	var builder strings.Builder

	// Build divider
	var dividerBuilder strings.Builder
	dividerBuilder.WriteByte('+')
	for i := 0; i < width+2; i++ {
		dividerBuilder.WriteRune('-')
	}
	dividerBuilder.WriteString("+\n")
	divider := dividerBuilder.String()

	builder.WriteString(divider)
	for _, entry := range entries {
		builder.WriteString(AddSideBorders(entry, width))
		builder.WriteString(divider)
	}

	return builder.String()
}

const (
	leftSideBorder  = "| %-"
	rightSideBorder = "s |"

	leftSideBorderless  = "  %-"
	rightSideBorderless = "s"
)

func AddSideBorders(in string, width int) string {
	var builder strings.Builder

	for _, line := range newlineRegexp.Split(in, -1) {
		if len(line) == 0 {
			builder.WriteString(fmt.Sprintf(
				leftSideBorder+strconv.Itoa(width)+rightSideBorder,
				"",
			))
			builder.WriteRune('\n')
		}

		for lIdx := 0; lIdx < len(line); lIdx += width {
			rIdx := lIdx + width
			leftSide := leftSideBorder
			rightSide := rightSideBorder

			if lIdx != 0 {
				leftSide = leftSideBorderless
			}

			if rIdx < len(line) {
				rightSide = rightSideBorderless
			} else {
				rIdx = len(line)
			}

			builder.WriteString(fmt.Sprintf(
				leftSide+strconv.Itoa(width)+rightSide,
				line[lIdx:rIdx],
			))
			builder.WriteRune('\n')
		}
	}

	return builder.String()
}
