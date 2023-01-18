package util

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strconv"
)

func NewSQ(db *sqlx.DB) sq.StatementBuilderType {
	return sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		RunWith(db)
}

func CheckErrNoRows(err error, message400, message500 string) error {
	if err == sql.ErrNoRows {
		return fiber.NewError(http.StatusBadRequest, message400)
	}
	return NewError(err, http.StatusInternalServerError, message500)
}

func NewError(err error, code int, message string) error {
	log.Println(err)
	return fiber.NewError(code, message)
}

func QueryUint64(ctx *fiber.Ctx, key string, defaultValue ...string) (n uint64) {
	n, _ = strconv.ParseUint(ctx.Query(key, defaultValue...), 10, 64)
	return n
}
