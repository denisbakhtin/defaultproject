package system

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/zenazn/goji/web"

	"github.com/gorilla/sessions"

	"html/template"

	"io"
	"os"
	"path/filepath"
	"strings"
)

type Application struct {
	Configuration *Configuration
	Template      *template.Template
	Store         *sessions.CookieStore
	DB            *sqlx.DB //it maintains a connection pool internally, thread safe
}

func (application *Application) Init(filename *string) {
	application.Configuration = &Configuration{}
	err := application.Configuration.Load(*filename)

	if err != nil {
		glog.Fatalf("Can't read configuration file: %s", err)
		panic(err)
	}

	application.Store = sessions.NewCookieStore([]byte(application.Configuration.Secret))
}

func (application *Application) LoadTemplates() error {
	var templates []string

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			templates = append(templates, path)
		}
		return nil
	}

	err := filepath.Walk(application.Configuration.TemplatePath, fn)

	if err != nil {
		return err
	}

	application.Template = template.Must(template.ParseFiles(templates...))
	return nil
}

func (application *Application) ConnectToDatabase() {
	var err error
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", application.Configuration.Database.Host, application.Configuration.Database.User, application.Configuration.Database.Password, application.Configuration.Database.Name)

	application.DB, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		glog.Fatalf("Can't connect to the database: %v", err)
		panic(err)
	}
}

func (application *Application) Close() {
	glog.Info("Bye!")
	application.DB.Close()
}

func (application *Application) Route(hand func(web.C, *http.Request) (string, int)) web.Handler {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		c.Env["Content-Type"] = "text/html"

		glog.Errorf("%+v\n", hand)
		body, code := hand(c, r)
		//body, code := "", 200

		if session, exists := c.Env["Session"]; exists {
			err := session.(*sessions.Session).Save(r, w)
			if err != nil {
				glog.Errorf("Can't save session: %v", err)
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
