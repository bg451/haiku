package main

import (
	"log"
	"net/http"
	"strconv"
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
