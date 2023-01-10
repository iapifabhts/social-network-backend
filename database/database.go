package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func New(config string) *sql.DB {
	db, err := sql.Open("postgres", config)
	if err != nil {
		panic(err)
	}
	return db
}
