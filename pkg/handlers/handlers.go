package handlers

import (
	"log"
	"net/http"

	"github.com/yesilyurtburak/go-web-basics-5/models"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/config"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/forms"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/render"
)

// Type definition for Repository pattern
type Repository struct {
	App *config.AppConfig
}

// Variable declaration for Repository pattern
var Repo *Repository

// Function definition for creating a new Repository
func NewRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}

// Function definition to handle routing with Repository pattern
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.gotmpl", &models.PageData{})
}

func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {
	strMap := make(map[string]string)
	render.RenderTemplate(w, r, "about.page.gotmpl", &models.PageData{StrMap: strMap})
	// created a strMap and send some information to about.page.gotmpl template via models.PageData
}

func (m *Repository) LoginHandler(w http.ResponseWriter, r *http.Request) {
	strMap := make(map[string]string)
	render.RenderTemplate(w, r, "login.page.gotmpl", &models.PageData{StrMap: strMap})
}

func (m *Repository) PageHandler(w http.ResponseWriter, r *http.Request) {
	strMap := make(map[string]string)
	render.RenderTemplate(w, r, "page.page.gotmpl", &models.PageData{StrMap: strMap})
}

func (m *Repository) MakePostHandler(w http.ResponseWriter, r *http.Request) {
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

	// get the form data from template and create a new struct `article`
	article := models.Article{
		BlogTitle:   r.Form.Get("blog_title"),
		BlogArticle: r.Form.Get("blog_article"),
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
