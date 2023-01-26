package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/config"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/handlers"
)

// in place of built-in packages, we use a 3rd party package for routing: chi

func routes(app *config.AppConfig) http.Handler {
	// create a multiplexer to handle incoming http requests
	mux := chi.NewMux()

	// using middlewares
	mux.Use(middleware.Recoverer) // built-in middleware for handling the panic gracefully.
	mux.Use(LogRequestInfo)       // custom middleware that is developed for a spesific purpose.
	mux.Use(NoSurf)
	mux.Use(SetupSession)

	// routing pages
	mux.Get("/", handlers.Repo.HomeHandler)
	mux.Get("/about", handlers.Repo.AboutHandler)
	mux.Get("/login", handlers.Repo.LoginHandler)
	mux.Post("/login", handlers.Repo.PostLoginHandler)
	mux.Get("/makepost", handlers.Repo.MakePostHandler)
	mux.Post("/makepost", handlers.Repo.PostMakePostHandler)
	mux.Get("/page", handlers.Repo.PageHandler)
	mux.Get("/article-received", handlers.Repo.ArticleReceived)

	// for static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
