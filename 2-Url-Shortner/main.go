package main

import (
	"2-Url-Shortner/urlshort"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := defaultMux()

	//Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	}
	var mapHandler = urlshort.MapHandler(pathsToUrls, mux)
	fileName := flag.String("file", "urls.db", "yaml/json/db file containing paths and urls")
	flag.Parse()

	var dataHandler http.HandlerFunc
	isDatabaseFile := strings.Contains(*fileName, "db")

	if isDatabaseFile {
		db, err := bolt.Open(*fileName, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		dataHandler = urlshort.DBHandler(db, mapHandler)

	} else {
		buf, err := os.ReadFile(*fileName)
		if err != nil {
			log.Fatal("Unable to read input file ", err)
		}
		dataHandler, err = urlshort.DataHandler(buf, mapHandler)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Starting the server on :8080")
	err := http.ListenAndServe(":8080", dataHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultPage)
	return mux
}

func defaultPage(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Unable to resolve path: %v\n", r.URL.Path)
	if err != nil {
		return
	}
}
