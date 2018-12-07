package main

import (
	"flag"
	"fmt"
	golog "log"
	"net/http"
	"os"
	"strconv"

	"github.com/andydennisonbooth/smock/echo"
	"github.com/andydennisonbooth/smock/log"
	"github.com/andydennisonbooth/smock/mock"
	"github.com/andydennisonbooth/smock/proxy"
)

func main() {
	port := flag.Int("p", 8889, "port to run server on")
	width := flag.Int("w", log.DefaultWidth, "width of output")

	echoFlag := flag.Bool("e", false, "run echo server")
	response := flag.String("r", "", "mock response")
	filename := flag.String("f", "", "filename of file containing mock response")
	proxyTarget := flag.String("x", "", "base url for proxy server")

	flag.Parse()

	mockServer := *filename != "" || *response != ""
	proxyServer := *proxyTarget != ""
	echoServer := *echoFlag

	var handler http.Handler
	var serverType string
	switch {
	case proxyServer:
		serverType = "Proxy"

		p, err := proxy.New(*proxyTarget, nil)
		if err != nil {
			golog.Fatal("Unable to make proxy server", err)
		}

		p.Width = *width
		handler = p

	case mockServer:
		serverType = "Mock"

		if *response != "" && *filename != "" {
			golog.Fatal("Too many responses provided")
		}

		var m *mock.Mock
		if *response != "" {
			m = mock.NewFromString(*response)
		} else if *filename != "" {
			var err error
			m, err = mock.NewFromFile(*filename)
			if err != nil {
				golog.Fatal("Unable to read file", err)
			}
		}
		m.Width = *width
		handler = m

	case echoServer:
		serverType = "Echo"

		e := echo.New()
		e.Width = *width
		handler = e

	default:
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Printf("%s server running @ localhost:%d\n", serverType, *port)
	golog.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), handler))
}
