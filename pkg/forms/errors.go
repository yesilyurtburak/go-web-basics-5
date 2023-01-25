package forms

type errors map[string][]string

// This function adds a new error to the errors map for given id.
func (e errors) AddError(tagID, message string) {
	e[tagID] = append(e[tagID], message)
}

// This function gets the error associated with given id, if exists.
func (e errors) GetError(tagID string) string {
	es := e[tagID]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
