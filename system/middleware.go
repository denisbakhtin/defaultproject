package system

import (
	"net/http"

	"github.com/denisbakhtin/defaultproject/models"
	"github.com/golang/glog"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji/web"
)

// Makes sure templates are stored in the context
func (application *Application) ApplyTemplates(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.Env["Template"] = application.Template
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Makes sure controllers can have access to session
func (application *Application) ApplySessions(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session, _ := application.Store.Get(r, "session")
		c.Env["Session"] = session
		h.ServeHTTP(w, r)
		context.Clear(r)
	}
	return http.HandlerFunc(fn)
}

// Makes sure controllers can have access to the database
func (application *Application) ApplyDatabase(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.Env["DB"] = application.DB
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (application *Application) ApplyAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session := c.Env["Session"].(*sessions.Session)
		if userId, ok := session.Values["UserId"].(int64); ok {
			database := c.Env["DB"].(*sqlx.DB)

			user, err := models.GetUser(database, userId)
			if err != nil {
				glog.Warningf("Auth error: %v", err)
				c.Env["User"] = nil
			} else {
				c.Env["User"] = user
			}
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
