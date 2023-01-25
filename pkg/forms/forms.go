package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// holds information that is assigned to our forms.
type Form struct {
	url.Values
	Errors errors
}

// This function creates a new form in our codebase.
func NewForm(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// This function loops through all fields and checks if they're empty or not.
func (f *Form) HasRequired(tagIDs ...string) {
	for _, tagID := range tagIDs {
		value := f.Get(tagID)
		if strings.TrimSpace(value) == "" {
			f.Errors.AddError(tagID, "This field cannot be blank!")
		}
	}
}

// This function checks if a field with given id has content in the form.
func (f *Form) HasValue(tagID string, r *http.Request) bool {
	x := r.Form.Get(tagID)
	return x != "" // returns true if field is not empty; returns false otherwise.
	// if x == "" {
	// 	f.Errors.AddError(tagID, "Field empty") // add error if has no content
	// 	return false
	// }
	// return true
}

// This function checks if the given field meets the minimum length condition.
func (f *Form) MinLength(tagID string, length int, r *http.Request) bool {
	x := r.Form.Get(tagID)
	if len(x) < length { // add an error if conditions are not met
		f.Errors.AddError(tagID, fmt.Sprintf("This field must be %d characters long or more.", length))
		return false
	}
	return true
}

// This function checks if the form is valid.
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0 // form is valid if there are no errors.
}

// ------------------------------------------------------------------------- //
// We can use 3rd party package for much wider checks. (asaskevich/govalidator)

// This function uses 3rd party package to check if the given email is valid or not.
func (f *Form) IsEmail(tagID string) {
	if !govalidator.IsEmail(f.Get(tagID)) {
		f.Errors.AddError(tagID, "Invalid email address!")
	}
}
