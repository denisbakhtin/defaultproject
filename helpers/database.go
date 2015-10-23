package helpers

import (
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji/web"
)

func GetDatabase(c web.C) *sqlx.DB {
	return c.Env["DB"].(*sqlx.DB)
}
