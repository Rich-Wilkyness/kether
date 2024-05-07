package forms

type errors map[string][]string // its a slice of []strings because we might have more than 1 error for a given field (ex: password; must have a capital, a number, and be 8 char long. if they only give abc, then we return 3 seperate errors)

// adds an error message to our type for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message) // field will be a field from our form and the message will be our error message we want them to see
}

// returns the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	return es[0]
}
