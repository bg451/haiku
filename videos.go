package main

import (
	"bytes"
	"fmt"
	"log"
)

type Video struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
	Elo int    `json:"elo"`
}

func (database *Database) getAll() string {
	var buffer bytes.Buffer

	rows, err := database.db.Query("SELECT * FROM videos")
	handleErr(err)

	for rows.Next() {
		var id int
		var url string
		var elo int

		err = rows.Scan(&id, &url, &elo)
		handleErr(err)
		buffer.WriteString(fmt.Sprintf("%d: %d, %s", id, elo, url))
	}
	return buffer.String()
}
func (database *Database) updateVideo(v *Video) {
	stmt := fmt.Sprintf("UPDATE videos SET elo=%d WHERE video_id=%d",
		v.Elo, v.ID)
	_, err := database.db.Exec(stmt)
	if err != nil {
		log.Printf("Error updating match result")
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
