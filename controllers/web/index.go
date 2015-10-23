package web

import (
	"net/http"

	"html/template"

	"github.com/denisbakhtin/defaultproject/helpers"
	"github.com/zenazn/goji/web"
)

// Home page route
func Index(c web.C, r *http.Request) (string, int) {
	t := helpers.GetTemplate(c)

	widgets := helpers.Parse(t, "home", nil)

	// With that kind of flags template can "figure out" what route is being rendered
	c.Env["IsIndex"] = true

	c.Env["Title"] = "Default Project - free Go website project template"
	c.Env["Content"] = template.HTML(widgets)

	return helpers.Parse(t, "main", c.Env), http.StatusOK
}
