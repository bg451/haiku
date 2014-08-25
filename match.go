package main

import (
	"fmt"
	"log"
)

type Match struct {
	ID      int
	video_a *Video
	video_b *Video
	winnerA bool
}

func createMatch(a *Video, b *Video, winnerA bool) *Match {
	//FIXME: Validate the inputs to make sure everything isnt nil
	return &Match{video_a: a, video_b: b, winnerA: winnerA}
}

//sql

func insertMatch(m *Match) {
	b := "FALSE"
	if m.winnerA {
		b = "TRUE"
	}
	sqlStmt := fmt.Sprintf(`INSERT INTO matches (video_a_id, video_b_id, winnerA)
                          VALUES (%d, %d, %s)`,
		m.video_a.ID,
		m.video_b.ID,
		b)
	_, err := database.db.Exec(sqlStmt)
	if err != nil {
		log.Printf("Could not execute statement: %q", err)
	}
}

func findMatchById(id int) (Match, error) {
	var (
		vidAId  int
		vidBId  int
		winnerA bool
	)
	sqlStmt := fmt.Sprintf("SELECT vid_a_id, vid_b_id, winnerA FROM matches WHERE match_id=%d", id)
	err := database.db.QueryRow(sqlStmt).Scan(&vidAId, &vidBId, &winnerA)
	if err != nil {
		return Match{}, fmt.Errorf("Could not find match by id")
	}
	vidA, err := findVideoById(vidAId)
	vidB, err := findVideoById(vidBId)
	if err != nil {
		return Match{}, err
	}
	return Match{ID: id, video_a: vidA, video_b: vidB, winnerA: winnerA}, nil
}
