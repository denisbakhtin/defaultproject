package main

import (
	"flag"
	"github.com/golang/glog"
	"net/http"

	"github.com/elcct/defaultproject/controllers"
	"github.com/elcct/defaultproject/system"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func main() {
	filename := flag.String("config", "config.json", "Path to configuration file")

	flag.Parse()
	defer glog.Flush()

	var application = &system.Application{}

	application.Init(filename)
	application.LoadTemplates()
	application.ConnectToDatabase()

	// Setup static files
	static := web.New()
	static.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(application.Configuration.PublicPath))))

	http.Handle("/assets/", static)

	// Apply middleware
	goji.Use(application.ApplyTemplates)
	goji.Use(application.ApplySessions)
	goji.Use(application.ApplyDatabase)
	goji.Use(application.ApplyAuth)

	controller := &controllers.MainController{}

	goji.Get("/signin", application.Route(controller, "SignIn"))
	goji.Post("/signin", application.Route(controller, "SignInPost"))

	goji.Get("/signup", application.Route(controller, "SignUp"))
	goji.Post("/signup", application.Route(controller, "SignUpPost"))

	goji.Get("/logout", application.Route(controller, "Logout"))

	goji.Get("/", application.Route(controller, "Index"))

	graceful.PostHook(func() {
		application.Close()
	})
	goji.Serve()
}