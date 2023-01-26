package render

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/yesilyurtburak/go-web-basics-5/models"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/config"
	"github.com/yesilyurtburak/go-web-basics-5/pkg/helpers"
)

// do not render everything while reloading the page, use templateCache instead.
var templateCache = make(map[string]*template.Template)

var app *config.AppConfig

func NewAppConfig(a *config.AppConfig) {
	app = a
}

// This function pass the CSRF token to the template to handle post requests.
func AddCSRFData(pd *models.PageData, r *http.Request) *models.PageData {
	pd.CSRFToken = nosurf.Token(r) // generate a token

	// add information to the pagedata if user logged in or not.
	if app.Session.Exists(r.Context(), "user_id") {
		pd.IsAuthenticated = true
	}

	return pd
}

// this function creates a cache.
func makeTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.gotmpl",
	}
	// this creates a new template from the given files
	tmpl, err := template.ParseFiles(templates...)
	helpers.ErrorCheck(err)
	templateCache[t] = tmpl // add the loaded page's template data to the templateCache map.
	return nil
}

// this function renders the template on the browser
func RenderTemplate(w http.ResponseWriter, r *http.Request, t string, pd *models.PageData) {
	var tmpl *template.Template
	var err error
	// check if the template is already in cache
	_, inMap := templateCache[t]
	if !inMap {
		err = makeTemplateCache(t)
		helpers.ErrorCheck(err)
	} else {
		fmt.Println("Template is loaded from cache")
	}
	tmpl = templateCache[t]

	pd = AddCSRFData(pd, r) // Add CSRFToken to the existing pd variable.

	err = tmpl.Execute(w, pd) // writes the template to the response writer `w` by sending data `pd`
	helpers.ErrorCheck(err)
}
