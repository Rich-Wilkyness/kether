package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Rich-Wilkyness/kether/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// IsAuthenticated checks if a user is authenticated
// by checking if the user_id session exists
func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user")
	return exists
}
