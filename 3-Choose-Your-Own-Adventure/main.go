package main

import (
	"3-Choose-Your-Own-Adventure/parser"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl *template.Template

func main() {
	const FileName = "story.json"

	buf, err := os.ReadFile(FileName)
	if err != nil {
		log.Fatal(err, FileName)
	}

	data, err := parser.ParseData(buf)
	if err != nil {
		log.Fatal(err)
	}

	tmpl = template.Must(template.ParseFiles("templates/index.gohtml"))
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, data[r.URL.Path[1:]])
		if err != nil {
			log.Fatal(err)
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
