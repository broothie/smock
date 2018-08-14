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

	filename := flag.String("f", "", "File to read response from for mock server")
	proxyTarget := flag.String("x", "", "Base url for proxy server")
	echo := flag.Bool("e", false, "Run echo server")

	flag.Parse()
	arg := flag.Arg(0)

	mockServer := *filename != "" || arg != ""
	proxyServer := *proxyTarget != ""
	echoServer := *echo

	if (mockServer && proxyServer) || (mockServer && echoServer) || (proxyServer && echoServer) {
		golog.Fatal("Only one type of server can be used at a time (between mock-, proxy-, and echo-servers.")
	}

	var handler http.HandlerFunc
	var serverType string
	switch {
	case proxyServer:
		serverType = "Proxy"
		handler = proxy.MustNewProxy(*proxyTarget, !*quiet, *width).ServeHTTP
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
		serverType = "Mock"

		if arg == "" && *filename == "" {
			golog.Fatal("No response provided")
		}

		if arg != "" && *filename != "" {
			golog.Fatal("Too many responses provided")
		}

		var response []byte
		if arg != "" {
			response = []byte(arg)
		}

		if *filename != "" {
			data, err := ioutil.ReadFile(*filename)
			if err != nil {
				golog.Fatal(err)
			}
			response = data
		}

		handler = func(w http.ResponseWriter, _ *http.Request) {
			w.Write(response)
		}
		if !*quiet {
			handler = log.Middleware(*width, handler)
		}
	}

	http.HandleFunc("/", handler)
	fmt.Printf("%s server running @ localhost:%s\n", serverType, *port)
	golog.Fatal(http.ListenAndServe(":"+*port, nil))
}
