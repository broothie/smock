package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andydennisonbooth/reqo"
)

const defaultPort = "8889"

func main() {
	// Parse args
	port := flag.String("p", defaultPort, "Port to run echo server on")
	responseFilename := flag.String("f", "", "File to read response from")
	flag.Parse()
	arg := flag.Arg(0)

	if arg == "" && *responseFilename == "" {
		log.Fatal("No response provided")
	}

	if arg != "" && *responseFilename != "" {
		log.Fatal("Too many responses provided")
	}

	var response []byte
	if arg != "" {
		response = []byte(arg)
	}

	if *responseFilename != "" {
		data, err := ioutil.ReadFile(*responseFilename)
		if err != nil {
			log.Fatal(err)
		}
		response = data
	}

	http.HandleFunc("/", reqo.LogMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write(response)
	}))
	fmt.Printf("Mockserver @ localhost:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
