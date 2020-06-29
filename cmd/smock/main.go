package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/broothie/smock/ui"

	"github.com/alecthomas/kingpin"
	"github.com/broothie/smock/handlers"
)

var (
	version, date string

	// Top level flags
	port   = kingpin.Flag("port", "port to run server mock on").Short('p').Default("9090").Int()
	uiPort = kingpin.Flag("uiport", "port to run ui on").Short('u').Default("9091").Int()
	noUI   = kingpin.Flag("no-ui", "disable ui").Short('U').Default("false").Bool()

	_ = kingpin.Command("version", "print smock version")

	// smock [mock]
	mock        = kingpin.Command("mock", "mock response").Default()
	mockCode    = mock.Flag("code", "response status code").Short('c').Default("200").Int()
	mockHeaders = mock.Flag("header", "response headers").Short('h').StringMap()
	mockBody    = mock.Flag("body", "response body, provided as a string").Short('b').Default("").String()
	mockFile    = mock.Flag("file", "response body, read from a file").Short('f').Default("").String()
	//mockStubFile = mock.Flag("stub", "stub file").Short('s').Default("").String()

	// smock proxy
	proxy    = kingpin.Command("proxy", "reverse proxy to target url").Alias("p")
	proxyURL = proxy.Arg("url", "url to proxy to").Required().URL()
)

func main() {
	logger := log.New(os.Stdout, "[smock] ", 0)

	var intro string
	var handler http.Handler
	switch kingpin.Parse() {
	case "mock":
		intro = fmt.Sprintf("mock server @ http://localhost:%d", *port)
		if *mockFile != "" {
			handler = handlers.File(*mockCode, stringMapToHTTPHeader(*mockHeaders), *mockFile)
			//} else if *mockStubFile != "" {
			//	var err error
			//	handler, err = stub.Handler(*mockStubFile)
			//	if err != nil {
			//		fmt.Println(err)
			//		os.Exit(1)
			//		return
			//	}
		} else {
			handler = handlers.Mock(logger, *mockCode, stringMapToHTTPHeader(*mockHeaders), []byte(*mockBody))
		}

	case "proxy":
		intro = fmt.Sprintf("proxying http://localhost:%d ➜ %v", *port, *proxyURL)
		handler = handlers.Proxy(*proxyURL)

	case "version":
		fmt.Printf("smock v%s; built %s\n", version, date)
		os.Exit(0)
		return

	default:
		fmt.Println("invalid command")
		os.Exit(1)
		return
	}

	uiServer := ui.New(logger, *uiPort)
	if !*noUI {
		go uiServer.Serve()
	}

	handler = uiServer.Middleware(handler)
	handler = dateSuppressor(handler)
	logger.Println(intro)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), handler); err != nil {
		logger.Println(err)
		os.Exit(1)
	}
}

func dateSuppressor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header()["Date"] = nil
		next.ServeHTTP(w, r)
	})
}

func stringMapToHTTPHeader(stringMap map[string]string) http.Header {
	header := make(http.Header)
	for key, value := range stringMap {
		header.Add(key, value)
	}

	return header
}
