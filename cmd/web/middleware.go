package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
)

// Middleware performs some action either before or after a request.

// This middleware function prints time and url path information to the terminal for each incoming request.
func LogRequestInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		fmt.Printf("%d/%d/%d : %d/%d ", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute())
		fmt.Println(r.URL.Path)
		next.ServeHTTP(w, r) // move on to the next data that we want to serve
	})
}

// A middleware for setup session data.
func SetupSession(next http.Handler) http.Handler {
	return sessionManager.LoadAndSave(next) // Load and save the session data by passing a session token as a cookie.
}

// A middleware for setup csrf protection.
func NoSurf(next http.Handler) http.Handler {
	noSurfHandler := nosurf.New(next)
	noSurfHandler.SetBaseCookie(http.Cookie{
		Name:     "mycsrfcookie",
		Path:     "/",
		Domain:   "",
		Secure:   false, // are cookies encrypted? true requires https
		HttpOnly: true,
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode, // for development
	})
	return noSurfHandler
}
