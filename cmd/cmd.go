package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	golog "log"
	"net/http"
	"net/http/httputil"

	"github.com/andydennisonbooth/smock/log"
	"github.com/andydennisonbooth/smock/proxy"
)

const widthOffset = 2

func main() {
	port := flag.String("p", "8889", "Port to run mock server on")
	quiet := flag.Bool("q", false, "Be quiet")
	width := flag.Int("w", 80, "Width of output")

	response := flag.String("r", "", "String to respond with")
	filename := flag.String("f", "", "File to read response from for mock server")
	proxyTarget := flag.String("x", "", "Base url for proxy server")
	echo := flag.Bool("e", false, "Run echo server")

	flag.Parse()

	mockServer := *filename != "" || *response != ""
	proxyServer := *proxyTarget != ""
	echoServer := *echo

	var handler http.HandlerFunc
	var serverType string
	switch {
	case proxyServer:
		serverType = "Proxy"
		handler = proxy.MustNewProxy(*proxyTarget, !*quiet, *width).ServeHTTP

	case mockServer:
		serverType = "Mock"

		if *response != "" && *filename != "" {
			golog.Fatal("Too many responses provided")
		}

		var res []byte
		if *response != "" {
			res = []byte(*response)
		}

		if *filename != "" {
			data, err := ioutil.ReadFile(*filename)
			if err != nil {
				golog.Fatal(err)
			}
			res = data
		}

		handler = func(w http.ResponseWriter, _ *http.Request) {
			w.Write(res)
		}
		if !*quiet {
			handler = log.Middleware(*width, handler)
		}

	case echoServer:
		serverType = "Echo"

		handler = func(w http.ResponseWriter, r *http.Request) {
			dump, _ := httputil.DumpRequest(r, r.Method != http.MethodGet)
			w.Write(dump)
		}

		if !*quiet {
			handler = log.Middleware(*width, handler)
		}

	default:
		golog.Fatal("For some reason, I don't know what to do")
	}

	http.HandleFunc("/", handler)
	fmt.Printf("%s server running @ localhost:%s\n", serverType, *port)
	golog.Fatal(http.ListenAndServe(":"+*port, nil))
}
