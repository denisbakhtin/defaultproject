package helpers

import (
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
)

func GetSession(c web.C) *sessions.Session {
	return c.Env["Session"].(*sessions.Session)
}
