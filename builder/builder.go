package builder

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

func New(config string) squirrel.StatementBuilderType {
	db, err := sql.Open("postgres", config)
	if err != nil {
		panic(err)
	}
	return squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(db)
}
