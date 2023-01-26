package helpers

import (
	"log"
	"net/http"

	"github.com/yesilyurtburak/go-web-basics-5/pkg/config"
)

var app *config.AppConfig

func ErrorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// check if user_id exists in the session
func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}
