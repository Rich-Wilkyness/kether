package models

import "github.com/Rich-Wilkyness/kether/internal/forms"

// -------------------------------------- Setup type of Data sent to Render ------------------
// holds whatever kind of data we think will be sent to be rendered in the template
type TemplateData struct {
	StringMap map[string]string // rather than using just a string type, we use a map of string so if we have 1 or 20 strings being passed to our render we can do it
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // this allows us to pass any kind of object to our render (interface{} means any kind of type that is an object)
	CSRFToken string                 // security token for forms. a Cross Site Request Forgery Token, needed on every page, because they might have a form of some kind
	// hidden field in forms that is a long string of numbers. handled usually with middleware (using nosurf)
	Flash   string      // this is a "flash message" success that something happened for the user to know (message sent)
	Warning string      // type of message to warn a user (page not saved, etc.)
	Error   string      // type of message to notify a user an error occured
	Form    *forms.Form // we put this here so we have access to our form throughout our app
}
