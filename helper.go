package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func handleErr(e error) {
	if e != nil {
		log.Printf("Error: %q\n", e)
	}
}
func parseInt(str string) (int, error) {
	i, err := strconv.ParseInt(str, 0, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}
func setJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func boolToInt(b bool) int {
	if b == true {
		return 1
	}
	return 0
}
func intToBool(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

func logRequest(r *http.Request) {
	log.Println("%s: %15s %s ", r.Method, r.URL.String(), r.RemoteAddr)
}

// Ping youtube to make sure the video exists
// Then validate the video length
// If the video is less than 35 seconds, create a new one

func validateUrl(urlString string) (string, error) {
	// oauthKey = "TABzAdrAofgB9Vw7NIffXgSl"
	// client, err := buildOAuthHTTPClient(oauthKey)
	// if err != nil {
	// 	return
	// }
	// service, err := youtube.New(client)
	// if err != nil {
	// 	return
	// }
	url, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("Invalid url: Cannot parse")
	}
	q := url.Query()
	host := url.Host
	if !strings.Contains(host, "youtube.com") {
		return "", fmt.Errorf("Invalid url: Host not youtube")
	}
	id := q.Get("v")
	// call := service.VideoListCall.Id(id)
	// results, err := call.Do()
	// if err != nil {
	// 	log.Printf("swag")
	// 	return
	// }
	return "//" + host + "/embed/" + id, nil
}
