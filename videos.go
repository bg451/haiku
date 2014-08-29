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

func (database *Database) getAll() {
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
func (database *Database) getRandomVideo() (*Video, error) {
	log.Printf("getRandomVideo")
	var id int
	var url string
	var elo int
	stmt := fmt.Sprintf("SELECT * FROM videos ORDER BY RANDOM() LIMIT 1")
	err := database.db.QueryRow(stmt).Scan(&id, &url, &elo)
	if err != nil {
		log.Printf("swag: %s", err.Error())
		return &Video{}, err
	}
	return &Video{id, url, elo}, err
}
func (database *Database) findVideoById(id int) (*Video, error) {
	log.Printf("findVideoById")
	var vId int
	var url string
	var elo int
	stmt := fmt.Sprintf("SELECT * FROM videos WHERE video_id=%d", id)
	err := database.db.QueryRow(stmt).Scan(&vId, &url, &elo)
	if err != nil {
		return &Video{}, fmt.Errorf("Video not found")
	}

	return &Video{vId, url, elo}, nil
}
func (database *Database) insertNewVideo(v Video) error {
	log.Printf("Inserting new video")
	v.Elo = 1500
	sqlStmt := fmt.Sprintf("INSERT INTO videos (elo, url) VALUES (%d, '%s')", v.Elo, v.Url)
	_, err := database.db.Exec(sqlStmt)
	handleErr(err)
	return err
}
