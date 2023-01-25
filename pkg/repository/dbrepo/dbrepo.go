package dbrepo

import (
	"database/sql"

	"github.com/yesilyurtburak/go-web-basics-5/pkg/config"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: app,
		DB:  conn,
	}
}
