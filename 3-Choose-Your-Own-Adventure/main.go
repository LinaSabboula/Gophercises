package main

import (
	"3-Choose-Your-Own-Adventure/parser"
	"html/template"
	"log"
	"net/http"
	"os"
)

var tmpl *template.Template

type Router struct {
	data map[string]parser.Chapter
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		r.URL.Path = "/intro"
	}
	err := tmpl.Execute(w, router.data[r.URL.Path[1:]])
	if err != nil {
		log.Fatal(err)
	}
}

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

	router := Router{data: data}

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
