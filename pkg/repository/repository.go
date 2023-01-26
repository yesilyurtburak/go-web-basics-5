package repository

import "github.com/yesilyurtburak/go-web-basics-5/models"

// allows us to list all of database functions that we wanna be able to access by all of our handlers.
type DatabaseRepo interface {
	InsertPost(newPost models.Post) error
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	AuthenticateUser(email, testPassword string) (int, string, error)
	AddUser(user models.User) error
	GetAnArticle() (int, int, string, string, error)
	GetArticlesForHomepage() (models.ArticleList, error)
}
