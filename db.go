package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func initDb(path string) *Database {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	verifyVideoTable(db)
	verifyMatchTable(db)
	return &Database{db: db}
}
func verifyVideoTable(db *sql.DB) {
	_, err := db.Exec("SELECT * FROM videos")
	if err != nil {
		log.Printf("Error: %q", err)
		log.Printf("Creating videos table")
		sqlStmt := `
		CREATE TABLE videos (video_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		                     url VARCHAR(100),
												 elo INTEGER);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("Could not create videos table: %q\n", err)
		}
	}
}

func verifyMatchTable(db *sql.DB) {
	_, err := db.Exec("SELECT * FROM matches")
	if err != nil {
		log.Printf("Error: %q", err)
		log.Printf("Creating matches table")
		sqlStmt := `
		CREATE TABLE matches (match_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
													video_a_id INTEGER,
													video_b_id INTEGER,
													winnerA INTEGER NOT NULL,
		                      FOREIGN KEY(video_a_id) REFERENCES videos(video_id),
													FOREIGN KEY(video_b_id) REFERENCES videos(video_id));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("Could not create matches table: %q\n", err)
		}
	}
}
