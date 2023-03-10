package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/yesilyurtburak/go-web-basics-5/models"
	"golang.org/x/crypto/bcrypt"
)

// Functions for accessing database

func (m *postgresDBRepo) InsertPost(newPost models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // cancels if can't connect within 5 seconds.

	query := `INSERT INTO post(user_id, title, content) VALUES ($1, $2, $3)`

	// execute a query with a context
	_, err := m.DB.ExecContext(ctx, query, newPost.UserID, newPost.Title, newPost.Content)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // cancels if can't connect within 5 seconds.

	query := `SELECT id, name, email, password, account_created, last_login, user_type FROM users WHERE id = $1`

	// execute a query with a context and returns at most 1 row.
	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.AccountCreated, &u.LastLogin, &u.UserType)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // cancels if can't connect within 5 seconds.

	query := `UPDATE users SET name = $1, email = $2, last_login = $3, user_type = $4`

	_, err := m.DB.ExecContext(ctx, query, u.Name, u.Email, time.Now(), u.UserType)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) AuthenticateUser(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // cancels if can't connect within 5 seconds.

	var id int
	var hashedPW string

	query := `SELECT id, password FROM users WHERE email=$1`

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&id, &hashedPW)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("password is incorrect")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPW, nil
}

func (m *postgresDBRepo) AddUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // cancels if can't connect within 5 seconds.

	// hash the user password before inserting into the database
	pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}

	hashedPassword := string(pw)

	query := `INSERT INTO users(name, email, password, account_created, last_login, user_type) VALUES ($1, $2, $3, $4, $5, $6);`

	_, err = m.DB.ExecContext(ctx, query, user.Name, user.Email, hashedPassword, user.AccountCreated, user.LastLogin, user.UserType)
	if err != nil {
		return err
	}
	return nil
}

// select an article from post table and fetch it.
func (m *postgresDBRepo) GetAnArticle() (int, int, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // cancels if can't connect within 5 seconds.

	var id, uID int
	var aTitle, aContent string

	query := `SELECT id, user_id, title, content FROM post LIMIT 1`

	row := m.DB.QueryRowContext(ctx, query)
	err := row.Scan(&id, &uID, &aTitle, &aContent)
	if err != nil {
		return id, uID, "", "", err
	}
	return id, uID, aTitle, aContent, nil
}

func (m *postgresDBRepo) GetArticlesForHomepage() (models.ArticleList, error) {
	var arList models.ArticleList
	query := `SELECT id, user_id, title, content FROM post ORDER BY id DESC LIMIT $1`
	rows, err := m.DB.Query(query, 3)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, uID int
		var title, content string
		err = rows.Scan(&id, &uID, &title, &content)
		if err != nil {
			panic(err)
		}
		arList.ID = append(arList.ID, id)
		arList.UserID = append(arList.UserID, uID)
		arList.Title = append(arList.Title, title)
		arList.Content = append(arList.Content, content)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return arList, nil
}
