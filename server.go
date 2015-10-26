package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/denisbakhtin/defaultproject/controllers/web"
	"github.com/denisbakhtin/defaultproject/system"

	"github.com/GeertJohan/go.rice"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	gojiweb "github.com/zenazn/goji/web"
)

func main() {
	env := flag.String("e", "development", "Application environment: development, production, testing")
	migration := flag.String("migrate", "", "Run DB migrations: up, down, redo, new [MIGRATION_NAME] and then os.Exit(0)")

	flag.Parse()
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.InfoLevel)

	var application = &system.Application{}

	//create rice boxes for folders with data
	configBox := rice.MustFindBox("config")         //config dir
	migrationsBox := rice.MustFindBox("migrations") //migrations dir
	viewsBox := rice.MustFindBox("views")           //views dir
	publicBox := rice.MustFindBox("public")         //public dir
	imagesBox := rice.MustFindBox("public/images")  //public/images dir

	application.Init(env, configBox)
	application.ConnectToDatabase()
	if len(*migration) > 0 {
		//Read https://github.com/rubenv/sql-migrate for more info about migrations
		application.RunMigrations(migrationsBox, migration)
		application.Close()
		os.Exit(0)
	}
	err := application.LoadTemplates(viewsBox)
	if err != nil {
		logrus.Fatal(err)
	}

	// Setup static files
	static := gojiweb.New()

	//static.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(application.Configuration.PublicPath))))
	static.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(publicBox.HTTPBox())))
	http.Handle("/assets/", static)

	// Couple of files - in the real world you would use nginx to serve them.
	goji.Get("/robots.txt", http.FileServer(publicBox.HTTPBox()))
	goji.Get("/favicon.ico", http.FileServer(imagesBox.HTTPBox()))

	//Apply middlewares
	goji.Use(application.ApplyTemplates)
	goji.Use(application.ApplySessions)
	goji.Use(application.ApplyDatabase)
	goji.Use(application.ApplyAuth)

	// Home page
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
