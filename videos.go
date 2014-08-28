package main

import (
	"fmt"
	"log"
)

type Video struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
	Elo int    `json:"elo"`
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
func getRandomVideo() (*Video, error) {
	var id int
	var url string
	var elo int
	stmt := "SELECT * FROM videos WHERE video_id >= (abs(random()) % (SELECT max(video_id) FROM videos) + 1) LIMIT 1;"
	row, err := database.db.Query(stmt)
	if err != nil {
		return &Video{}, err
	}
	row.Next()
	err = row.Scan(&id, &url, &elo)
	handleErr(err)
	return &Video{id, url, elo}, err

}
func findVideoById(id int) (*Video, error) {
	var vId int
	var url string
	var elo int
	stmt := fmt.Sprintf("SELECT * FROM videos WHERE video_id=%d", id)
	row, err := database.db.Query(stmt)
	if err != nil {
		return &Video{}, fmt.Errorf("Video not found")
	}
	row.Next()
	err = row.Scan(&vId, &url, &elo)
	if err != nil {
		return &Video{}, fmt.Errorf("%q", err)
	}

	return &Video{vId, url, elo}, nil
}
func insertNewVideo(v Video) error {
	v.Elo = 1500
	sqlStmt := fmt.Sprintf("INSERT INTO videos (elo, url) VALUES (%d, '%s')", v.Elo, v.Url)
	_, err := database.db.Exec(sqlStmt)
	handleErr(err)
	return err
}
func (v *Video) update() {
	_, err := database.db.Exec("UPDATE videos SET elo=%d WHERE id=%d", v.Elo, v.ID)
	handleErr(err)
}
