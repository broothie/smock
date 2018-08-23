package main

import (
	"flag"
	"fmt"
	golog "log"
	"net/http"

	"github.com/andydennisonbooth/smock/echo"
	"github.com/andydennisonbooth/smock/log"
	"github.com/andydennisonbooth/smock/mock"
	"github.com/andydennisonbooth/smock/proxy"
)

func main() {
	port := flag.String("p", "8889", "Port to run mock server on")
	width := flag.Int("w", log.DefaultWidth, "Width of output")

	echoFlag := flag.Bool("e", false, "Run echo server")
	response := flag.String("r", "", "String to respond with")
	filename := flag.String("f", "", "File to read response from for mock server")
	proxyTarget := flag.String("x", "", "Base url for proxy server")

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
		fallthrough
	default:
		serverType = "Echo"

		e := echo.New()
		e.Width = *width
		handler = e
	}

	fmt.Printf("%s server running @ localhost:%s\n", serverType, *port)
	golog.Fatal(http.ListenAndServe(":"+*port, handler))
}
