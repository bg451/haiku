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

func (database *Database) getVideosSorted() string {
	var buffer bytes.Buffer

	rows, err := database.db.Query("SELECT * FROM videos ORDER BY elo DESC")
	handleErr(err)

	for rows.Next() {
		var id int
		var url string
		var elo int

		err = rows.Scan(&id, &url, &elo)
		handleErr(err)
		buffer.WriteString(fmt.Sprintf("%d: %d, %s\n", id, elo, url))
	}
	return buffer.String()
}
func (database *Database) updateVideo(v *Video) {
	log.Printf("\t\t\tStarting video update")
	stmt := fmt.Sprintf("UPDATE videos SET elo=%d WHERE video_id=%d",
		v.Elo, v.ID)
	_, err := database.db.Exec(stmt)
	if err != nil {
		log.Printf("Error updating match result")
	}
	log.Printf("\t\t\tEnding video update")
}
func (database *Database) getRandomVideo(excludeA int, excludeB int) (*Video, error) {
	log.Printf("\t\tStarting getRandomVideo")
	var id int
	var url string
	var elo int
	stmt := fmt.Sprintf("SELECT * FROM videos OFFSET random() * (select count(*) from videos) limit 1")
	err := database.db.QueryRow(stmt).Scan(&id, &url, &elo)
	if id == excludeA || id == excludeB {
		err = database.db.QueryRow(stmt).Scan(&id, &url, &elo)
	}
	log.Printf("\t\tEnding getRandomVideo")
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
func (database *Database) findVideoInELoRange(elo int, id int) (*Video, error) {
	log.Printf("\t\t Starting findVideoInELoRange")
	var vId int
	var url string
	var swag int
	//Start off with a +- 10 Elo difference, expand by 5 until found
	uR := elo + 20
	lR := elo - 20
	stmt := fmt.Sprintf(`
  SELECT * FROM videos WHERE (elo BETWEEN %d and %d)
  AND NOT video_id=%d ORDER BY RANDOM() LIMIT 1;`, lR, uR, id)
	err := database.db.QueryRow(stmt).Scan(&vId, &url, &swag)

	for {
		if err != nil {
			uR = uR + 10
			lR = lR - 10
			if lR <= 0 {
				err = fmt.Errorf("Video not found: ELO range exceeded 0")
				break
			}
			stmt := fmt.Sprintf(`
      SELECT * FROM videos WHERE (elo BETWEEN %d and %d)
      AND NOT video_id=%d ORDER BY RANDOM() LIMIT 1;`, lR, uR, id)
			err = database.db.QueryRow(stmt).Scan(&vId, &url, &swag)
		} else {
			break
		}
	}
	log.Printf("\t\t Ending findVideoInELoRange")
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
