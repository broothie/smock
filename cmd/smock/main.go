package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/broothie/smock/pkg/handlers"
	"github.com/broothie/smock/pkg/reqlogger"
	"github.com/broothie/smock/pkg/ui"
)

var (
	version, date string

	port   = kingpin.Flag("port", "port to run smock on").Short('p').Default("9090").Int()
	uiPort = kingpin.Flag("uiport", "port to run ui on").Short('u').Default("9091").Int()
	skipUI = kingpin.Flag("no-ui", "disable ui").Default("false").Bool()

	_ = kingpin.Command("version", "print smock version")

	// smock [mock]
	mock        = kingpin.Command("mock", "mock response based on command args").Default()
	mockCode    = mock.Flag("code", "status code").Short('c').Default("200").Int()
	mockHeaders = mock.Flag("header", "headers").Short('h').StringMap()
	mockBody    = mock.Flag("body", "body").Short('b').Default("").String()

	// smock file
	file     = kingpin.Command("file", "mock response from file")
	fileName = file.Arg("filename", "file to mock response with").Required().String()

	// smock proxy
	proxy       = kingpin.Command("proxy", "reverse proxy to target url")
	proxyTarget = proxy.Arg("target", "url to proxy to").Required().URL()
)

func main() {
	logger := log.New(os.Stdout, "[smock] ", 0)

	var intro string
	var handler http.Handler
	switch kingpin.Parse() {
	case "mock":
		intro = fmt.Sprintf("mock server @ http://localhost:%d", *port)
		handler = handlers.Mock(logger, *mockCode, stringMapToHTTPHeader(*mockHeaders), []byte(*mockBody))
	case "file":
		intro = fmt.Sprintf("responding with '%s' contents @ http://localhost:%d", *fileName, *port)
		handler = handlers.File(logger, *fileName)
	case "proxy":
		intro = fmt.Sprintf("proxying http://localhost:%d → %s", *port, *proxyTarget)
		handler = handlers.Proxy(*proxyTarget)
	case "version":
		fmt.Printf("fileserver v%s; built %s\n", version, date)
		os.Exit(0)
	default:
		log.Println("invalid command")
		os.Exit(1)
	}

	if *skipUI {
		handler = reqlogger.New(logger)(handler)
	} else {
		ui := ui.New(*uiPort, logger)
		handler = ui.Middleware(handler)
		ui.Start()
	}

	logger.Println(intro)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), handler); err != nil {
		logger.Println(err)
		os.Exit(1)
	}
}

func stringMapToHTTPHeader(stringMap map[string]string) http.Header {
	header := make(http.Header)
	for key, value := range stringMap {
		header.Set(key, value)
	}

	return header
}
