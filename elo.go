package main

import (
	"fmt"
	"math"
)

// Elo K constant as per wikipedia
const K = 30

type Video struct {
	id  int
	Url string
	Elo int
}

func expectedRating(a *Video, b *Video) (exA float64, exB float64) {
	// Calculate expected score based on wikipedia
	// http://en.wikipedia.org/wiki/Elo_rating_system#Mathematical_details
	// eA = 1 / (1 + 10^((eloB - eloA)/400))
	exA = 1.0 / (1.0 + math.Pow(10, float64(b.Elo-a.Elo)/400.0))
	exB = 1.0 / (1.0 + math.Pow(10, float64(a.Elo-b.Elo)/400.0))
	return exA, exB
}

func runMatch(a *Video, b *Video, winnerA bool) {
	var (
		sA float64 = 0.0
		sB float64 = 0.0
	)
	exA, exB := expectedRating(a, b)
	fmt.Printf("exA: %d, exB: %d", exA, exB)
	// Assuming A won, set the sB to 0 and sA to 1
	if winnerA == true {
		sA = 1
	} else {
		sB = 1
	}

	// Calculate the new ELO's
	// R' = R + K(score - expectedScore)
	a.Elo = a.Elo + int(math.Floor(K*(sA-exA)))
	b.Elo = b.Elo + int(math.Floor(K*(sB-exB)))
}
