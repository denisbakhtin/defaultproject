package system

import (
	"fmt"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/zenazn/goji/web"

	"github.com/gorilla/sessions"

	"html/template"

	"io"
	"os"
	"strings"
)

type Application struct {
	Configuration *Configuration
	Template      *template.Template
	Store         *sessions.CookieStore
	DB            *sqlx.DB //has internal threadsafe connection pool
}

func (application *Application) Init(env *string, box *rice.Box) {
	var err error
	application.Configuration, err = LoadConfiguration(env, box.MustBytes("config.json"))

	if err != nil {
		logrus.Fatalf("Can't read configuration file: %s", err)
		panic(err)
	}

	application.Store = sessions.NewCookieStore([]byte(application.Configuration.Secret))
}

func (application *Application) LoadTemplates(box *rice.Box) error {
	tmpl := template.New("")

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			var err error
			tmpl, err = tmpl.Parse(box.MustString(path))
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := box.Walk("", fn)
	if err != nil {
		return err
	}

	application.Template = tmpl
	return nil
}

func (application *Application) ConnectToDatabase() {
	var err error
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", application.Configuration.Database.Host, application.Configuration.Database.User, application.Configuration.Database.Password, application.Configuration.Database.Name)

	application.DB, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		logrus.Fatalf("Can't connect to the database: %v", err)
		panic(err)
	}
}

func (application *Application) RunMigrations(box *rice.Box, command *string) {
	switch *command {
	case "new":
		migrateNew(box)
	case "up":
		migrateUp(application.DB.DB, box)
	case "down":
		migrateDown(application.DB.DB, box)
	case "redo":
		migrateDown(application.DB.DB, box)
		migrateUp(application.DB.DB, box)
	default:
		logrus.Fatalf("Wrong migration flag %q, acceptable values: up, down", *command)
	}
}

func (application *Application) Close() {
	logrus.Info("Bye!")
	application.DB.Close()
}

func (application *Application) Route(action func(web.C, *http.Request) (string, int)) web.Handler {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		c.Env["Content-Type"] = "text/html"

		body, code := action(c, r)

		if session, exists := c.Env["Session"]; exists {
			err := session.(*sessions.Session).Save(r, w)
			if err != nil {
				logrus.Errorf("Can't save session: %v", err)
			}
		}

		switch code {
		case http.StatusOK:
			if _, exists := c.Env["Content-Type"]; exists {
				w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
			}
			io.WriteString(w, body)
		case http.StatusSeeOther, http.StatusFound:
			http.Redirect(w, r, body, code)
		default:
			w.WriteHeader(code)
			io.WriteString(w, body)
		}
	}
	return web.HandlerFunc(fn)
}
