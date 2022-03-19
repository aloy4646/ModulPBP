package controller

import (
	"database/sql"
	"log"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_latihan_pbp")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
