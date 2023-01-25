package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/yesilyurtburak/go-web-basics-5/models"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/config"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/dbdriver"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/handlers"
)

const portNumber = "8080"
const ipAddress = "127.0.0.1"

var url = fmt.Sprintf("%s:%s", ipAddress, portNumber)

// sessions: remember the information about users when they traverse our site. We track this information by saving a cookie on the user's device. (using a 3rd party for managing cookies : alexedwards/scs)
var sessionManager *scs.SessionManager // defines a new SessionManager variable `sessionManager`
var app config.AppConfig               // defines a new configuration variable `app`

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	// create and configure a new server
	srv := &http.Server{
		Addr:    url,
		Handler: routes(&app),
	}

	// listen to the traffic for incoming http requests.
	fmt.Printf("Listening traffic at %s\n", url)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*dbdriver.DB, error) {
	gob.Register(models.Article{}) // can use models.Article inside of our sessions.
	gob.Register(models.User{})    // adds User table model into the session
	gob.Register(models.Post{})    // adds Post table model into the session

	// initialize the session
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour              // cookies' lifetime
	sessionManager.Cookie.Persist = true                  // remain the cookies even if when the browser closes.
	sessionManager.Cookie.Secure = false                  // are cookies encrypted? false for Development mode; true requires HTTPS.
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode // Development:SameSiteLaxMode, Production:SameSiteStrictMode

	app.Session = sessionManager // saves the session information to the system.

	// create and connect to the database
	db, err := dbdriver.ConnectSQL("host=localhost port=5432 dbname=blog_db user=postgres password=postgres")
	if err != nil {
		log.Fatal("Can't connect to the database.")
	}

	repo := handlers.NewRepo(&app, db) // creates a new repo and a new database repo
	handlers.NewHandlers(repo)         // this assign the value `repo` to `Repo` variable inside the handlers.go

	return db, nil
}
