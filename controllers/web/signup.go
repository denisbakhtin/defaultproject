package web

import (
	"html/template"
	"net/http"

	"github.com/denisbakhtin/defaultproject/helpers"
	"github.com/denisbakhtin/defaultproject/models"
	"github.com/golang/glog"
	"github.com/zenazn/goji/web"
)

// Sign up route
func SignUp(c web.C, r *http.Request) (string, int) {
	t := helpers.GetTemplate(c)
	session := helpers.GetSession(c)

	// With that kind of flags template can "figure out" what route is being rendered
	c.Env["IsSignUp"] = true

	c.Env["Flash"] = session.Flashes("auth")

	var widgets = helpers.Parse(t, "auth/signup", c.Env)

	c.Env["Title"] = "Default Project - Sign Up"
	c.Env["Content"] = template.HTML(widgets)

	return helpers.Parse(t, "main", c.Env), http.StatusOK
}

// Sign Up form submit route. Registers new user or shows Sign Up route with appropriate messages set in session
func SignUpPost(c web.C, r *http.Request) (string, int) {
	email, password := r.FormValue("email"), r.FormValue("password")

	session := helpers.GetSession(c)
	database := helpers.GetDatabase(c)

	user, err := models.GetUserByEmail(database, email)

	if user != nil {
		session.AddFlash("User exists", "auth")
		return SignUp(c, r)
	}

	user = &models.User{
		Name:  email,
		Email: email,
	}
	err = user.HashPassword(password)
	if err != nil {
		session.AddFlash("Error whilst registering user.", "auth")
		glog.Errorf("Error whilst registering user: %v", err)
		return SignUp(c, r)
	}

	if err := models.InsertUser(database, user); err != nil {
		session.AddFlash("Error whilst registering user.", "auth")
		glog.Errorf("Error whilst registering user: %v", err)
		return SignUp(c, r)
	}

	session.Values["UserId"] = user.Id

	return "/", http.StatusSeeOther
}
