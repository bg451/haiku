package main

import (
	"fmt"
	"testing"
)

func TestExpectedRating(t *testing.T) {
	vA := &Video{ID: 1, Url: "swag", Elo: 2100}
	vB := &Video{ID: 2, Url: "beast", Elo: 2000}
	exA, exB := expectedRating(vA, vB)
	sExA := fmt.Sprintf("%1.11f", exA)
	sExB := fmt.Sprintf("%1.7f", exB)
	if sExA != "0.64006499980" {
		t.Error(sExA)
	}
	if sExB != "0.3599350" {
		t.Error(sExB)
	}
}
func TestRunMatch_aWin(t *testing.T) {
	vA := &Video{ID: 1, Url: "swag", Elo: 2100}
	vB := &Video{ID: 2, Url: "beast", Elo: 2000}
	calculateElo(vA, vB, true)
	if vA.Elo != 2110 {
		t.Error(fmt.Sprintf("Expected Elo: 2110, Got: %d", vA.Elo))
	}
	if vB.Elo != 1989 {
		t.Error(fmt.Sprintf("Expected Elo: 1989, Got: %d", vB.Elo))
	}
}

func TestRunMatch_bWin(t *testing.T) {
	vA := &Video{ID: 1, Url: "swag", Elo: 2100}
	vB := &Video{ID: 2, Url: "beast", Elo: 2000}
	calculateElo(vA, vB, false)
	if vA.Elo != 2080 {
		t.Error(fmt.Sprintf("Expected Elo: 2110, Got: %d", vA.Elo))
	}
	if vB.Elo != 2019 {
		t.Error(fmt.Sprintf("Expected Elo: 1989, Got: %d", vB.Elo))
	}
}
