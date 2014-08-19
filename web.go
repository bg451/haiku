package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

const httpPort = 8080

var db *Database

type Post struct {
	Swag string
}

func main() {
	database = new(Database)
	database.db = initDb("./db.sqlite3")
	go webServer()
	for {
	}
}
func webServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", indexPageHandler).Methods("GET")
	r.HandleFunc("/{blog}", indexPageHandler).Methods("GET")
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static/"))))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("templates/blog.html", "templates/base.html")
	if err != nil {
		fmt.Fprintf(w, "There was a an error %s", err)
		return
	}
	templ.ExecuteTemplate(w, "base", nil)
}
