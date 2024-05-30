package main

import (
	"net/http"

	"github.com/Rich-Wilkyness/kether/internal/config"
	"github.com/Rich-Wilkyness/kether/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// this library simplifies routing as we will see and provides middleware authentication
func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer) // the documentation on this is on the github
	mux.Use(NoSurf)               // this is our middleware for csrf toxen authentication
	mux.Use(SessionLoad)          // this helps us have state management

	mux.Get("/", handlers.Repo.Home)

	mux.Get("/make-test", handlers.Repo.Test)
	mux.Post("/make-test", handlers.Repo.PostTest)

	mux.Get("/user/register", handlers.Repo.Register)
	mux.Post("/user/register", handlers.Repo.PostRegister)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)

	// this allows our tmpl templates to access our static directory
	// this directory is where we will store things like images
	fileServer := http.FileServer(http.Dir("./static/"))             // we first find our directory. by using "./" this is our root. and this is what is required by the Dir function
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // then we feed our mux the directory by directing it to our static directory and removing static from the pathname of our files to get our filename

	// this protects our admin routes
	// the routes will have a combine name to get to them. for example /admin/dashboard
	mux.Route("/account", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.ShowAccountDashboard)
		mux.Get("/edit", handlers.Repo.ShowAccountEdit)
		mux.Post("/edit", handlers.Repo.PostAccountEdit)
		mux.Post("/delete", handlers.Repo.PostAccountDelete)
	})

	return mux
}
