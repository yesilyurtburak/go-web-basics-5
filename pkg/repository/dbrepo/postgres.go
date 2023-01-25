package dbrepo

import (
	"context"
	"time"

	"github.com/yesilyurtburak/go-web-basics-5/models"
)

// Functions for accessing database

func (m *postgresDBRepo) InsertPost(newPost models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // cancels if can't connect within 5 seconds.

	query := `INSERT INTO posts(user_id, title, content) VALUES ($1, $2, $3)`

	// execute a query with a context
	_, err := m.DB.ExecContext(ctx, query, newPost.ID, newPost.Title, newPost.Content)
	if err != nil {
		return err
	}
	return nil
}
