package main

import (
	"Gophercises/2-Url-Shortner/urlshort"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := defaultMux()

	//Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	}
	var mapHandler = urlshort.MapHandler(pathsToUrls, mux)
	fileName := flag.String("file", "urls.yml", "yaml file containing paths and urls")
	flag.Parse()

	buf, err := os.ReadFile(*fileName)

	if err != nil {
		log.Fatal("Unable to read input file ", err)
	}

	yamlHandler, err := urlshort.DataHandler(buf, mapHandler)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshort.DataHandler(buf, yamlHandler)
	fmt.Println("Starting the server on :8080")
	err = http.ListenAndServe(":8080", jsonHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
