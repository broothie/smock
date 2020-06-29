package stub

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func Handler(fileOrDir string) (http.HandlerFunc, error) {
	stubs, err := Parse(fileOrDir)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for _, stub := range stubs {
			if len(stub.Request.Methods) != 0 {
				if !stringInclude(stub.Request.Methods, r.Method) {
					continue
				}
			}

			if stub.Request.Path != "" {
				re, err := regexp.Compile(stub.Request.Path)
				if err != nil {
					continue
				}

				if !re.MatchString(r.URL.Path) {
					continue
				}

				captures := re.FindAllStringSubmatch(r.URL.Path, -1)[0:]
				fmt.Println(captures)
			}

			// Match
			for key, value := range stub.Response.Headers {
				w.Header().Add(key, value)
			}

			if stub.Response.Code != 0 {
				w.WriteHeader(stub.Response.Code)
			} else {
				w.WriteHeader(0)
			}

			w.Write([]byte(stub.Response.Body))
			return
		}

		w.Header().Add("x-smock-stub", "miss")
		w.WriteHeader(http.StatusTeapot)
	}, nil
}

func stringInclude(slice []string, includes string) bool {
	for _, str := range slice {
		if strings.EqualFold(str, includes) {
			return true
		}
	}

	return false
}
