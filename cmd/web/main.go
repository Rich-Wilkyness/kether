package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Rich-Wilkyness/kether/internal/config"
	"github.com/Rich-Wilkyness/kether/internal/driver"
	"github.com/Rich-Wilkyness/kether/internal/handlers"
	"github.com/Rich-Wilkyness/kether/internal/models"
	"github.com/Rich-Wilkyness/kether/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig // We want our global template cache here.
// this is global here to allow us to use the global variables in our middleware and other places.
// if it was inside main we couldn't call it in our middleware file

var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close() // this is a good practice to close the database connection when the program ends (run function ends)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app), // we point to our app
	}

	// start a web server that is listening
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// this allows us to store our type struct in a session. the session library by default can do strings and ints,
	// but not specific objects or types unless we do this (where we are registering the type to be stored)
	// we can store Question.TestID in a session for easy access to questions for a specific test
	gob.Register(models.Test{})
	gob.Register(models.User{})
	gob.Register(models.Question{})
	gob.Register(models.AnswerChoices{})
	gob.Register(models.Class{})
	gob.Register(models.TopicSelection{})
	gob.Register(models.UserInClass{})

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour              // We are assigning the lifetime of the session/cookie to 24 hours
	session.Cookie.Persist = true                  // when true, this means that after someone closes the browser/webpage the session ends
	session.Cookie.SameSite = http.SameSiteLaxMode // how strick is the site the cookie applies to. laxmode is default
	session.Cookie.Secure = app.InProduction       // for production we want this to be true, when true this makes the site https. and the cookies are encrypted

	app.Session = session

	// connect to the database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=postgres user=postgres password=lebrum1203")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc // add the template sets to the global cache
	app.UseCache = false

	// setup handlers
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	// this gives our render package access to the AppConfig (cache of template sets)
	render.NewRender(&app)

	return db, nil
}
