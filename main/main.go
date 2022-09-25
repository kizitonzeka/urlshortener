package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/kizitonzeka/urlshortener"
)

func main() {
	mux := defaultMux()

	urls := flag.String("url-file", "urls.yaml", " use url-file to parse shortened paths and their urls")

	flag.Parse()

	content, err := ioutil.ReadFile(*urls)

	if err != nil {
		log.Printf("Error reading file: %v", err)
		return
	}

	var urlHandler http.Handler

	if strings.Split(*urls, ".")[1] == "yaml" {
		log.Printf("Reading file: %v", *urls)
		urlHandler, err = urlshortener.YAMLHandler(content, mux)
		if err != nil {
			panic(err)
		}
	} else if strings.Split(*urls, ".")[1] == "json" {
		log.Printf("Reading file: %v", *urls)
		urlHandler, err = urlshortener.JSONHandler(content, mux)
		if err != nil {
			panic(err)
		}
	} else {
		log.Fatalf("Error reading file: %s. Please use a .yaml or .json file", *urls)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", urlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
