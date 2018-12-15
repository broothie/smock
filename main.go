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
	response := flag.String("r", "", "run mock server with provided response")
	filename := flag.String("f", "", "run mock server using file contents as response")
	proxyTarget := flag.String("x", "", "run proxy server pointed at provided uri")

	flag.Parse()

	mockServer := *filename != "" || *response != ""
	proxyServer := *proxyTarget != ""
	echoServer := *echoFlag

	var handler http.Handler
	var serverDetails string
	switch {
	case proxyServer:
		serverDetails = fmt.Sprintf("Reverse proxy to %q", *proxyTarget)

		p, err := proxy.New(*proxyTarget, nil)
		if err != nil {
			golog.Fatal("unable to make proxy server", err)
		}

		p.Width = *width
		handler = p

	case mockServer:
		serverDetails = "Mock server"

		if *response != "" && *filename != "" {
			golog.Fatal("too many responses provided")
		}

		var m *mock.Mock
		if *response != "" {
			m = mock.NewFromString(*response)
		} else if *filename != "" {
			var err error
			m, err = mock.NewFromFile(*filename)
			if err != nil {
				golog.Fatal("unable to read file", err)
			}
		}
		m.Width = *width
		handler = m

	case echoServer:
		serverDetails = "Echo server"

		e := echo.New()
		e.Width = *width
		handler = e

	default:
		flag.PrintDefaults()
		os.Exit(0)
	}

	fmt.Printf("%s running @ localhost:%d\n", serverDetails, *port)
	golog.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), handler))
}
