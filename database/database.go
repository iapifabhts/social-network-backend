package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New() *sqlx.DB {
	db, err := sqlx.Open("postgres",
		"postgresql://postgres:4100@localhost:4100/social_network?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}
