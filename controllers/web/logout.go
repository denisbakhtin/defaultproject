package web

import (
	"net/http"

	"github.com/denisbakhtin/defaultproject/helpers"
	"github.com/zenazn/goji/web"
)

// This route logs user out
func Logout(c web.C, r *http.Request) (string, int) {
	session := helpers.GetSession(c)

	delete(session.Values, "UserId")

	return "/", http.StatusSeeOther
}
