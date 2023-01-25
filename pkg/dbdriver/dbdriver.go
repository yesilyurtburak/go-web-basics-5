package dbdriver

/*
I have used a 3rd party package for managing database.
	$go get -u jackc/pgx/v5
*/

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// this holds our database connection
type DB struct {
	SQL *sql.DB
}

// initialize a reference to our DB type
var dbConn = &DB{}

// define rules for database
const maxOpenConnections = 20
const maxIdleConnections = 10
const maxLifetime = 5 * time.Minute

// This function creates a database
func NewDatabase(ds string) (*sql.DB, error) {
	db, err := sql.Open("pgx", ds) // connect to the database
	if err != nil {
		return nil, err
	}
	err = db.Ping() // ping the database
	if err != nil {
		return nil, err
	}
	return db, nil
}

// This function creates connection pool
func ConnectSQL(ds string) (*DB, error) {
	db, err := NewDatabase(ds)
	if err != nil {
		panic(err) // panics if you can't connect to the database
	}

	// set database rules
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(maxLifetime)

	dbConn.SQL = db
	return dbConn, nil
}
