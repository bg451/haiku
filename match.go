package main

import (
	"fmt"
	"log"
)

type Match struct {
	ID        int    `json:"id"`
	Video_a   *Video `json:"video_a"`
	Video_b   *Video `json:"video_b"`
	WinnerA   bool   `json:"winnerA"`
	Committed bool
}

func createMatch(a *Video, b *Video) *Match {
	return &Match{Video_a: a, Video_b: b}
}

func (database *Database) runMatch(m *Match) error {
	wA := boolToInt(m.WinnerA)
	calculateElo(m.Video_a, m.Video_b, m.WinnerA)
	database.updateVideo(m.Video_a)
	database.updateVideo(m.Video_b)

	stmt := fmt.Sprintf("UPDATE matches SET committed=%d, winnerA=%d WHERE match_id=%d",
		1, wA, m.ID)
	_, err := database.db.Exec(stmt)
	if err != nil {
		log.Printf("Error updating match result")
		return err
	}
	return nil
}
func (database *Database) generateMatch() (*Match, error) {
	log.Printf("Generate Match")
	vA, err := database.getRandomVideo()
	handleErr(err)
	vB, err := database.getRandomVideo()
	handleErr(err)
	for {
		if vB.ID != vA.ID {
			break
		}
		vB, err = database.getRandomVideo()
		if err != nil {
			log.Printf("Loop error %s", err.Error())
		}
	}
	match := &Match{Video_a: vA, Video_b: vB, WinnerA: false, Committed: false}
	dbase.insertMatch(match)
	return match, err
}

//sql

func (database *Database) insertMatch(m *Match) {
	wA := boolToInt(m.WinnerA)
	co := boolToInt(m.Committed)
	stmt := fmt.Sprintf("INSERT INTO matches (video_a_id, video_b_id, winnerA, committed) VALUES (%d, %d, %d,%d)",
		m.Video_a.ID, m.Video_b.ID, wA, co)
	result, err := database.db.Exec(stmt)
	if err != nil {
		log.Fatal(err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ID err: %s", err.Error())
	}
	m.ID = int(id)
}

func (database *Database) findMatchById(id int) (Match, error) {
	log.Printf("findmatchbyId")
	var (
		vidAId    int
		vidBId    int
		winnerA   int
		committed int
	)
	sqlStmt := fmt.Sprintf("SELECT * FROM matches WHERE match_id=%d", id)
	err := database.db.QueryRow(sqlStmt).Scan(&id, &vidAId, &vidBId, &winnerA, &committed)
	if err != nil {
		log.Printf(err.Error())
		return Match{}, fmt.Errorf("Could not find match by id")
	}
	vidA, err := database.findVideoById(vidAId)
	handleErr(err)
	vidB, err := database.findVideoById(vidBId)
	handleErr(err)
	if err != nil {
		return Match{}, err
	}
	return Match{ID: id, Video_a: vidA, Video_b: vidB, WinnerA: intToBool(winnerA), Committed: intToBool(committed)}, nil
}
