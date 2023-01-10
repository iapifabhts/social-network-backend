package util

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strconv"
)

func CheckErrNoRows(err error, message400, message500 string) error {
	log.Println(err)
	if err == sql.ErrNoRows {
		return fiber.NewError(http.StatusBadRequest, message400)
	}
	return fiber.NewError(http.StatusInternalServerError, message500)
}

func NewError(code int, message string, err error) error {
	log.Println(err)
	return fiber.NewError(code, message)
}

func SubQuery(selectBuilder squirrel.SelectBuilder) string {
	subQuery, _, _ := selectBuilder.ToSql()
	return fmt.Sprintf("(%s)", subQuery)
}

func QueryUint64(ctx *fiber.Ctx, key, defaultValue string) (n uint64) {
	n, _ = strconv.ParseUint(ctx.Query(key, defaultValue), 10, 64)
	return n
}

func NewPSQL(db *sql.DB) func() squirrel.StatementBuilderType {
	return func() squirrel.StatementBuilderType {
		return squirrel.StatementBuilder.
			PlaceholderFormat(squirrel.Dollar).
			RunWith(db)
	}
}
