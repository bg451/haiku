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

func (database *Database) getVideosSorted() (videos []Video) {

	rows, err := database.db.Query("SELECT video_id, url, elo FROM videos ORDER BY elo DESC")
	handleErr(err)

	for rows.Next() {
		var id int
		var url string
		var elo int

		err = rows.Scan(&id, &url, &elo)
		handleErr(err)
		videos = append(videos, Video{id, url, elo})
	}
	return videos
}
func (database *Database) updateVideo(v *Video) {
	stmt := fmt.Sprintf("UPDATE videos SET elo=%d WHERE video_id=%d",
		v.Elo, v.ID)
	_, err := database.db.Exec(stmt)
	if err != nil {
		log.Printf("Error updating match result")
	}
}
func (database *Database) getRandomVideo(excludeA int, excludeB int) (*Video, error) {
	var id int
	var url string
	var elo int
	stmt := fmt.Sprintf("SELECT video_id, url, elo FROM videos OFFSET random() * (select count(*) from videos) limit 1")
	err := database.db.QueryRow(stmt).Scan(&id, &url, &elo)
	if id == excludeA || id == excludeB {
		err = database.db.QueryRow(stmt).Scan(&id, &url, &elo)
	}
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
	stmt := fmt.Sprintf("SELECT video_id, url, elo FROM videos WHERE video_id=%d", id)
	err := database.db.QueryRow(stmt).Scan(&vId, &url, &elo)
	if err != nil {
		return &Video{}, fmt.Errorf("Video not found")
	}

	return &Video{vId, url, elo}, nil
}
func (database *Database) findVideoInELoRange(elo int, id int) (*Video, error) {
	var vId int
	var url string
	var swag int
	//Start off with a +- 10 Elo difference, expand by 5 until found
	uR := elo + 30
	lR := elo - 30
	stmt := fmt.Sprintf(`
  SELECT video_id, url, elo FROM videos WHERE (elo BETWEEN %d and %d)
  AND NOT video_id=%d ORDER BY RANDOM() LIMIT 1;`, lR, uR, id)
	err := database.db.QueryRow(stmt).Scan(&vId, &url, &swag)

	for {
		if err != nil {
			uR = uR + 30
			lR = lR - 30
			if lR <= 0 {
				err = fmt.Errorf("Video not found: ELO range exceeded 0")
				break
			}
			stmt := fmt.Sprintf(`
      SELECT video_id, url, elo FROM videos WHERE (elo BETWEEN %d and %d)
      AND NOT video_id=%d ORDER BY RANDOM() LIMIT 1;`, lR, uR, id)
			err = database.db.QueryRow(stmt).Scan(&vId, &url, &swag)
		} else {
			break
		}
	}
	return &Video{vId, url, swag}, err
}
func (database *Database) insertNewVideo(v Video) error {
	log.Printf("Inserting new video")
	v.Elo = 1500
	sqlStmt := fmt.Sprintf("INSERT INTO videos (elo, url) VALUES (%d, '%s')", v.Elo, v.Url)
	_, err := database.db.Exec(sqlStmt)
	handleErr(err)
	return err
}
