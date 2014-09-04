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

func (database *Database) runMatch(m *Match) {
	wA := boolToInt(m.WinnerA)
	calculateElo(m.Video_a, m.Video_b, m.WinnerA)
	database.updateVideo(m.Video_a)
	database.updateVideo(m.Video_b)

	stmt := fmt.Sprintf("UPDATE matches SET committed=%d, winnerA=%d WHERE match_id=%d",
		1, wA, m.ID)
	_, err := database.db.Exec(stmt)
	if err != nil {
		log.Printf("Error updating match result")
	}
}
func (database *Database) generateMatch(excludeA int, excludeB int) (*Match, error) {
	log.Printf("\t Starting GenerateMatch")
	vA, err := database.getRandomVideo(excludeA, excludeB)
	if err != nil {
		return &Match{}, err
	}
	vB, err := database.findVideoInELoRange(vA.Elo, vA.ID)
	if err != nil {
		return &Match{}, err
	}
	match := &Match{Video_a: vA, Video_b: vB, WinnerA: false, Committed: false}
	dbase.insertMatch(match)
	log.Printf("\tEnding GenerateMatch")
	return match, err
}

//sql

func (database *Database) insertMatch(m *Match) {
	var id int
	log.Printf("\t\tStrting matchInsertion")
	wA := boolToInt(m.WinnerA)
	co := boolToInt(m.Committed)
	stmt := fmt.Sprintf("INSERT INTO matches (video_a_id, video_b_id, winnerA, committed) VALUES (%d, %d, %d,%d) RETURNING match_id",
		m.Video_a.ID, m.Video_b.ID, wA, co)
	err := database.db.QueryRow(stmt).Scan(&id)
	if err != nil {
		log.Printf(err.Error())
	}
	if err != nil {
		log.Printf("ID err: %s", err.Error())
	}
	m.ID = int(id)
	log.Printf("\t\tEnding matchInsertion")
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
