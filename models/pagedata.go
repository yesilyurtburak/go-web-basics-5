package models

import "github.com/yesilyurtburak/go-web-basics-5/pkg/forms"

// These types can be used for sending data from our codebase to page templates.

type PageData struct {
	StrMap          map[string]string
	IntMap          map[string]int
	FltMap          map[string]float32
	DataMap         map[string]interface{} // value can be any type
	CSRFToken       string
	Error           string
	Warning         string
	Form            *forms.Form
	IsAuthenticated bool // user is currently logged in or not
}
