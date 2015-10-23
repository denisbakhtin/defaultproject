package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"

	"github.com/denisbakhtin/defaultproject/controllers/web"
	"github.com/denisbakhtin/defaultproject/system"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	gojiweb "github.com/zenazn/goji/web"
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
	static := gojiweb.New()
	static.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(application.Configuration.PublicPath))))

	http.Handle("/assets/", static)

	// Apply middleware
	goji.Use(application.ApplyTemplates)
	goji.Use(application.ApplySessions)
	goji.Use(application.ApplyDatabase)
	goji.Use(application.ApplyAuth)

	// Couple of files - in the real world you would use nginx to serve them.
	goji.Get("/robots.txt", http.FileServer(http.Dir(application.Configuration.PublicPath)))
	goji.Get("/favicon.ico", http.FileServer(http.Dir(application.Configuration.PublicPath+"/images")))

	// Home page
	//goji.Get("/", application.Route(controller, "Index"))
	goji.Get("/", application.Route(web.Index))

	// Sign In routes
	goji.Get("/signin", application.Route(web.SignIn))
	goji.Post("/signin", application.Route(web.SignInPost))

	// Sign Up routes
	goji.Get("/signup", application.Route(web.SignUp))
	goji.Post("/signup", application.Route(web.SignUpPost))

	// KTHXBYE
	goji.Get("/logout", application.Route(web.Logout))

	graceful.PostHook(func() {
		application.Close()
	})
	goji.Serve()
}
