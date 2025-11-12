package main

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
	"github.com/google/uuid"
)


db, err := sql.Open("sqlite3", "./mydb.db")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

err := db.Ping()