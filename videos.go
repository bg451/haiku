package main

import (
	"fmt"
	"log"
)

type Video struct {
	ID  int
	Url string
	Elo int
}

func getAll() {
	rows, err := database.db.Query("SELECT * FROM videos")
	handleErr(err)

	for rows.Next() {
		var id int
		var url string
		var elo int

		err = rows.Scan(&id, &url, &elo)
		handleErr(err)
		log.Printf("%d: %d, %s", id, elo, url)
	}
}
func findVideoById(id int) (*Video, error) {
	return &Video{}, nil
}
func (v *Video) update() {
	_, err := database.db.Exec("UPDATE videos SET elo=%d WHERE id=%d", v.Elo, v.ID)
	handleErr(err)
}
func (v *Video) insert() {
	sqlStmt := fmt.Sprintf("INSERT INTO videos (elo, url) VALUES (%d, '%s')", v.Elo, v.Url)
	_, err := database.db.Exec(sqlStmt)
	handleErr(err)
}
