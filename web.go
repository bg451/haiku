package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

const httpPort = 8080

var dbase *Database = new(Database)

type httpApiFunc func(http.ResponseWriter, *http.Request)

func main() {
	var err error
	dbase, err = initDb()
	if err != nil {
		log.Fatal(err.Error())
	}
	go startServer()
	for {
		// You need the sleep to prevent the program from locking up
		time.Sleep(time.Minute * 5)
		resp, err := http.Get("http://haitube.herokuapp.com")
		handleErr(err)
		if resp.StatusCode == 200 {
			log.Printf("server pinged")
		}
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
			"/matches":             getMatches,
			"/matches/new":         getMatchesNew,
			"/matches/{id:[0-9]+}": getMatchesID,
		},
		"POST": {
			"/videos/new":     postVideosNew,
			"/matches/result": postMatchesResult,
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Port: " + port)
	http.Handle("/", newRouter())
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	log.Printf("Serving http server")
	http.ListenAndServe(":"+port, nil)
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
	log.Printf("Starting request getNewMatches")
	r.ParseForm()
	va, err := strconv.ParseInt(r.Form.Get("idA"), 0, 0)
	vb, err := strconv.ParseInt(r.Form.Get("idB"), 0, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	match, err := dbase.generateMatch(int(va), int(vb))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, _ := json.Marshal(match)
	setJson(w)
	fmt.Fprintf(w, string(resp))
	log.Printf("Ending request getNewMatches")
}

func getMatchesID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 0, 0)
	if err != nil {
		fmt.Fprintf(w, "Could not parse int")
		return
	}
	swag := int(id)

	match, err := dbase.findMatchById(swag)
	if err != nil {
		fmt.Fprintf(w, "%d Error: %q", id, err)
		return
	}
	fmt.Fprintf(w, "The winner is %t", match.WinnerA)
}
func getVideos(w http.ResponseWriter, r *http.Request) {
	videos := dbase.getVideosSorted()
	templ, err := template.ParseFiles("templates/leaderboard.html", "templates/base.html")
	handleErr(err)
	templ.ExecuteTemplate(w, "base", videos)

}

func getVideosRandom(w http.ResponseWriter, r *http.Request) {
	log.Printf("getVideosRandom")

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

	result, err := dbase.findVideoById(id)
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

func postVideosNew(w http.ResponseWriter, r *http.Request) {
	vid := Video{}
	data, err := ioutil.ReadAll(r.Body)
	handleErr(err)
	err = json.Unmarshal(data, &vid)
	handleErr(err)
	url, err := validateUrl(vid.Url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vid.Url = url
	err = dbase.insertNewVideo(vid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func postMatchesResult(w http.ResponseWriter, r *http.Request) {
	log.Printf("Starting postMatchesResult")
	m := Match{}
	data, err := ioutil.ReadAll(r.Body)
	handleErr(err)
	err = json.Unmarshal(data, &m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go dbase.runMatch(&m)

	w.WriteHeader(http.StatusCreated)
	log.Printf("Ending postMatchesResult")
}
