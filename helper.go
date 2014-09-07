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

// Turns a youtube link into an embed link
func validateUrl(urlString string) (string, error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("Invalid url: Cannot parse")
	}
	q := url.Query()
	host := url.Host
	if !strings.Contains(host, "youtube.com") {
		return "", fmt.Errorf("Invalid url: Host not youtube")
	}
	return "//" + host + "/embed/" + q.Get("v"), nil
}
