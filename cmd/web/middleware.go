package main

import (
	"net/http"

	"github.com/Rich-Wilkyness/kether/internal/helpers"
	"github.com/justinas/nosurf"
)

// adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",                  // this means it hits the entire site
		Secure:   app.InProduction,     // in production this will be true (https) the s for secure. our global config file will change this to true when in production
		SameSite: http.SameSiteLaxMode, // standard
	})
	return csrfHandler
}

// loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler { // this will just load the session. this is important for state management
	return session.LoadAndSave(next) // helps remember state essentially
}

// Auth checks if a user is authenticated
func Auth(next http.Handler) http.Handler {
	// not sure how we have access to the w or r here. guess it is passed from next.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
