package log

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
)

const DefaultWidth = 80

var tabRegexp = regexp.MustCompile(`\t`)
var newlineRegexp = regexp.MustCompile(`\r?\n`)

func Entrify(width int, entries ...string) string {
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
		builder.WriteString(Borderify(entry, width))
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

func Borderify(in string, width int) string {
	var builder strings.Builder

	in = tabRegexp.ReplaceAllString(in, "  ")

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

func CleanDump(r *http.Request) (string, error) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return "", err
	}

	contentLength := r.Header.Get("Content-Length")
	if contentLength == "" {
		return string(dump[:len(dump)-4]), nil
	}

	return string(dump), nil
}
