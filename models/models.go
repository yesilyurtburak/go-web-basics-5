package models

import "time"

// Mimic our database tables with structs

type User struct {
	ID             int
	Name           string
	Email          string
	Password       string
	AccountCreated time.Time
	LastLogin      time.Time
	UserType       int
}

type Post struct {
	ID      int
	Title   string
	Content string
	UserID  int
}
