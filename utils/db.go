package utils

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "db.sqlite3"
const initDBSchema = `CREATE TABLE "files" (
	"id"	INTEGER,
	"filename"	INTEGER NOT NULL UNIQUE,
	PRIMARY KEY("id")
);`

var DB *sql.DB

func InitDb() {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatal("Something went wrong creating the database file")
		}
		file.Close()

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(initDBSchema)

		if err != nil {
			log.Fatal("Something went wrong creating the database schema")
		}
		DB = db
		return

	} else if err != nil {
		log.Fatal("Something went wrong checking if the database file exists")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
