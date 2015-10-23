package web

import (
	"net/http"

	"html/template"

	"github.com/denisbakhtin/defaultproject/helpers"
	"github.com/zenazn/goji/web"
)

// Sign in route
func SignIn(c web.C, r *http.Request) (string, int) {
	t := helpers.GetTemplate(c)
	session := helpers.GetSession(c)

	// With that kind of flags template can "figure out" what route is being rendered
	c.Env["IsSignIn"] = true

	c.Env["Flash"] = session.Flashes("auth")
	var widgets = helpers.Parse(t, "auth/signin", c.Env)

	c.Env["Title"] = "Default Project - Sign In"
	c.Env["Content"] = template.HTML(widgets)

	return helpers.Parse(t, "main", c.Env), http.StatusOK
}

// Sign In form submit route. Logs user in or set appropriate message in session if login was not succesful
func SignInPost(c web.C, r *http.Request) (string, int) {
	email, password := r.FormValue("email"), r.FormValue("password")

	session := helpers.GetSession(c)
	database := helpers.GetDatabase(c)

	user, err := helpers.Login(database, email, password)

	if err != nil {
		session.AddFlash("Invalid Email or Password", "auth")
		return SignIn(c, r)
	}

	session.Values["UserId"] = user.Id

	return "/", http.StatusSeeOther
}
