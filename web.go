package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

const httpPort = 8080

var database *Database = new(Database)

type httpApiFunc func(http.ResponseWriter, *http.Request)

func main() {
	var err error
	database, err = initDb("./foo.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	go startServer()
	for {
		// You need the sleep to prevent the program from locking up
		time.Sleep(time.Millisecond * 100)
	}
}
func addHandlers(router *mux.Router) {
	m := map[string]map[string]httpApiFunc{
		"GET": {
			"/":                    getIndexPage,
			"/videos":              getVideos,
			"/videos/{id:[0-9]+}":  getVideosID,
			"/videos/rand":         getVideosRandom,
			"/videos/new":          getVideosNew,
			"/videos/elo/highest":  getVideosEloHighest,
			"/matches":             getMatches,
			"/matches/new":         getMatchesNew,
			"/matches/{id:[0-9]+}": getMatchesID,
		},
		"POST": {
			"/videos/new":    postVideosNew,
			"matches/result": postMatchesResult,
		},
	}
	for method, routes := range m {
		for route, fnc := range routes {
			router.HandleFunc(route, fnc).Methods(method)
		}
	}
}

func newRouter() http.Handler {
	router := mux.NewRouter()
	addHandlers(router)
	return router
}
func startServer() {
	http.Handle("/", newRouter())
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	log.Printf("Serving http server")
	http.ListenAndServe(":8080", nil)
}
func getIndexPage(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("templates/blog.html", "templates/base.html")
	if err != nil {
		fmt.Fprintf(w, "There was a an error %s", err.Error())
		return
	}
	templ.ExecuteTemplate(w, "base", nil)
}
func getMatches(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getting matches")
}
func getMatchesNew(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getting matches/new")

}
func getMatchesID(w http.ResponseWriter, r *http.Request) {
	log.Printf("asfasfasd")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 0, 0)
	if err != nil {
		fmt.Fprintf(w, "Could not parse int")
		return
	}
	swag := int(id)

	match, err := findMatchById(swag)
	if err != nil {
		fmt.Fprintf(w, "%d Error: %q", id, err)
		return
	}
	fmt.Fprintf(w, "The winner is %t", match.winnerA)
}
func getVideos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getting videos")

}
func getVideosRandom(w http.ResponseWriter, r *http.Request) {
	vid, err := getRandomVideo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(vid)
	handleErr(err)
	setJson(w)
	fmt.Fprintf(w, string(json))
}
func getVideosNew(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFiles("templates/new_video.html", "templates/base.html")
	handleErr(err)
	templ.ExecuteTemplate(w, "base", nil)
}

func getVideosID(w http.ResponseWriter, r *http.Request) {
	log.Printf("getting videos/:id")
	vars := mux.Vars(r)
	id, err := parseInt(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	result, err := findVideoById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("/videos/%d Error: %s", id, err.Error())
		return
	}
	res, err := json.Marshal(result)
	handleErr(err)
	setJson(w)
	fmt.Fprintf(w, string(res))

}
func getVideosEloHighest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getting videos/elo/highest")

}

func postVideosNew(w http.ResponseWriter, r *http.Request) {
	vid := Video{}
	data, err := ioutil.ReadAll(r.Body)
	handleErr(err)
	err = json.Unmarshal(data, &vid)
	handleErr(err)
	log.Printf("POST: %s", string(data))

	err = insertNewVideo(vid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func postMatchesResult(w http.ResponseWriter, r *http.Request) {

}
