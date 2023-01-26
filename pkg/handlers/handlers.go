package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/yesilyurtburak/go-web-basics-5/models"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/config"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/dbdriver"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/forms"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/render"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/repository"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/repository/dbrepo"
)

// Type definition for Repository pattern
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Variable declaration for Repository pattern
var Repo *Repository

// Function definition for creating a new Repository
func NewRepo(app *config.AppConfig, db *dbdriver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresRepo(db.SQL, app),
	}
}

// Function definition to handle routing with Repository pattern
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// is user logged in to see this page?
	if !m.App.Session.Exists(r.Context(), "user_id") {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}

	// id, uid, title, content, err := m.DB.GetAnArticle()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Println("ID:", id)
	// fmt.Println("USER ID:", uid)
	// fmt.Println("TITLE:", title)
	// fmt.Println("CONTENT:", content)

	var articleList models.ArticleList
	articleList, err := m.DB.GetArticlesForHomepage()
	if err != nil {
		log.Println(err)
		return
	}

	// it will pass the data to the template
	data := make(map[string]interface{})
	data["articleList"] = articleList

	render.RenderTemplate(w, r, "home.page.gotmpl", &models.PageData{DataMap: data})
}

func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	// is user logged in to see this page?
	if !m.App.Session.Exists(r.Context(), "user_id") {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
	strMap := make(map[string]string)
	render.RenderTemplate(w, r, "about.page.gotmpl", &models.PageData{StrMap: strMap})
	// created a strMap and send some information to about.page.gotmpl template via models.PageData
}

func (m *Repository) LoginHandler(w http.ResponseWriter, r *http.Request) {
	strMap := make(map[string]string)
	render.RenderTemplate(w, r, "login.page.gotmpl", &models.PageData{StrMap: strMap})
}

func (m *Repository) PageHandler(w http.ResponseWriter, r *http.Request) {
	// is user logged in to see this page?
	if !m.App.Session.Exists(r.Context(), "user_id") {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
	strMap := make(map[string]string)
	render.RenderTemplate(w, r, "page.page.gotmpl", &models.PageData{StrMap: strMap})
}

func (m *Repository) MakePostHandler(w http.ResponseWriter, r *http.Request) {

	// is user logged in to see this page?
	if !m.App.Session.Exists(r.Context(), "user_id") {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}

	// created an empty article which we will then populate using the information on the server side that the user provides to us.
	var emptyArticle models.Article
	data := make(map[string]interface{})
	data["article"] = emptyArticle

	render.RenderTemplate(w, r, "makepost.page.gotmpl", &models.PageData{
		Form:    forms.NewForm(nil),
		DataMap: data,
	})
}

func (m *Repository) PostMakePostHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm() // sets the r.Form and r.PostForm
	if err != nil {
		log.Fatal(err)
		return
	}

	userId := m.App.Session.Get(r.Context(), "user_id")
	// get the form data from template and create a new struct `article`
	article := models.Post{
		UserID:  userId.(int),
		Title:   r.Form.Get("blog_title"),
		Content: r.Form.Get("blog_article"),
	}

	form := forms.NewForm(r.PostForm) // creates a new post form
	// form.HasValue("blog_title", r)    // checks the field and add an error if it is empty
	// form.HasValue("blog_article", r)  // checks the field and add an error if it is empty
	form.HasRequired("blog_title", "blog_article") // can add as many field parameters as you want

	form.MinLength("blog_title", 5, r)    // min. length should be 5 or more for the blog_title field.
	form.MinLength("blog_article", 10, r) // min. length should be 10 or more for the blog_article field.
	// form.IsEmail("email")  // ------> this can be used later.

	// if form has errors rerender the page and show the errors on the fields.
	if !form.IsValid() {
		data := make(map[string]interface{})
		data["article"] = article
		render.RenderTemplate(w, r, "makepost.page.gotmpl", &models.PageData{Form: form, DataMap: data})
		return
	}

	// Write article to the DB
	err = m.DB.InsertPost(article)
	if err != nil {
		log.Fatal(err)
	}

	m.App.Session.Put(r.Context(), "article", article) // put the post data to the session as k:v pair.
	http.Redirect(w, r, "/article-received", http.StatusSeeOther)
}

func (m *Repository) ArticleReceived(w http.ResponseWriter, r *http.Request) {
	// because Session.Get() returns an interface, we should assert the returned type as models.Article
	// retrieve the article data from the session
	article, ok := m.App.Session.Get(r.Context(), "article").(models.Article)
	if !ok {
		log.Println("Can't get data from the session")

		// if can't get data from session -> redirect to home page temporarily.
		m.App.Session.Put(r.Context(), "error", "Can't get data from the session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}
	data := make(map[string]interface{})
	data["article"] = article // I will get this data on article-received template.

	render.RenderTemplate(w, r, "article-received.page.gotmpl", &models.PageData{DataMap: data})
}

// new login
func (m *Repository) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context()) // prevent session fixation attempts

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.NewForm(r.PostForm)
	form.HasRequired("email", "password")
	form.IsEmail("email")

	if !form.IsValid() {
		render.RenderTemplate(w, r, "login.page.gotmpl", &models.PageData{
			Form: form,
		})
		return
	}

	// check if the entered password is same as the password that in the database.
	id, _, err := m.DB.AuthenticateUser(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid email or password")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id) // user logged in with success
	m.App.Session.Put(r.Context(), "flash", "Valid Login")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// This function logs out and clear all session.
func (m *Repository) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())    // delete all session data
	_ = m.App.Session.RenewToken(r.Context()) // prevent session fixation attempts
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// register new user
func (m *Repository) SignupHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "signup.page.gotmpl", &models.PageData{})
}

// register new user
func (m *Repository) PostSignupHandler(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context()) // prevent session fixation attempts

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	form := forms.NewForm(r.PostForm)

	form.HasRequired("name", "email", "password")
	form.IsEmail("email")
	form.MinLength("password", 5, r)

	if !form.IsValid() {
		render.RenderTemplate(w, r, "signup.page.gotmpl", &models.PageData{
			Form: form,
		})
		return
	}

	newUser := models.User{
		Name:           r.Form.Get("name"),
		Password:       r.Form.Get("password"),
		Email:          r.Form.Get("email"),
		UserType:       1,
		AccountCreated: time.Now(),
		LastLogin:      time.Now(),
	}

	err = m.DB.AddUser(newUser)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Registration failed!")
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Signed Up Successfully")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
