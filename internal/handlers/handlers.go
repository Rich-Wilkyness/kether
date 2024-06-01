package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Rich-Wilkyness/kether/internal/config"
	"github.com/Rich-Wilkyness/kether/internal/driver"
	"github.com/Rich-Wilkyness/kether/internal/forms"
	"github.com/Rich-Wilkyness/kether/internal/models"
	"github.com/Rich-Wilkyness/kether/internal/render"
	"github.com/Rich-Wilkyness/kether/internal/repository"
	"github.com/Rich-Wilkyness/kether/internal/repository/dbrepo"
)

// --------------------------------------- Setup Repositories ----------------------------------
// repository used by the handlers
var Repo *Repository

// sets the type of repository
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo // this is the interface we created in the repository package allows us to use any database we want in all of our handlers - aka passing information to the handlers
}

// creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// --------------------------------------- Handlers ------------------------------------------------
// for web browsers to work we need a response and request
// (m *Repository) gives access to the handlers everything that is inside repository
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Test(w http.ResponseWriter, r *http.Request) {
	var emptyTest models.Test // this is an empty test, by creating this with the data, we can populate our form with information that they previously input
	data := make(map[string]interface{})
	data["test"] = emptyTest
	render.Template(w, r, "make-test.page.tmpl", &models.TemplateData{
		Form: forms.New(nil), // creates a new form without anything in it when the make-test page is loaded.
		Data: data,
	})
}

func (m *Repository) PostTest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	classID, err := strconv.Atoi(r.Form.Get("class_id")) // we need to convert the class_id (string) to an int. Atoi means alpha (like alphanumeric) to int
	if err != nil {
		log.Println(err)
		return
	}
	// this field will be populated by the session, so I'm not sure if it needs to be converted (what type is it?)
	userID, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		log.Println(err)
		return
	}

	test := models.Test{
		Name:    r.Form.Get("name"),
		Version: r.Form.Get("version"),
		ClassID: classID,
		UserID:  userID,
	}

	form := forms.New(r.PostForm)

	form.Required("name", "version", "class_id", "user_id") // this is our server side validation for the form
	// form.MinLength("first_name", 3, r)
	// form.IsEmail("email")

	// this is so we can repopulate the form on errors, so they do not have to redo the form
	if !form.Valid() {
		data := make(map[string]interface{})
		data["test"] = test

		render.Template(w, r, "make-test.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// we need to insert our reservation into the database
	err = m.DB.InsertTest(test)
	if err != nil {
		log.Println(err)
		return
	}

	// m.App.Session.Put(r.Context(), "test", test) // don't need this

	// we want to redirect so our user does not submit twice (this is crucial for financial posts)
	http.Redirect(w, r, "/make-question", http.StatusSeeOther) // we can use other status' for redirect, but this one works well enough (303)
}

func (m *Repository) Question(w http.ResponseWriter, r *http.Request) {
	var emptyQuestion models.Question
	var emptyAnswerChoices models.AnswerChoices
	data := make(map[string]interface{})
	data["question"] = emptyQuestion
	data["answerChoices"] = emptyAnswerChoices
	render.Template(w, r, "make-question.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostQuestion(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	// convert the form data (string) to an int
	// 1) convert question model fields
	difficultyLevel, err := strconv.Atoi(r.Form.Get("difficulty_level"))
	if err != nil {
		log.Println(err)
		return
	}
	testID, err := strconv.Atoi(r.Form.Get("test_id"))
	if err != nil {
		log.Println(err)
		return
	}
	userID, err := strconv.Atoi(r.Form.Get("user_id"))
	if err != nil {
		log.Println(err)
		return
	}
	topicSelectionID, err := strconv.Atoi(r.Form.Get("topic_selection_id"))
	if err != nil {
		log.Println(err)
		return
	}

	question := models.Question{
		Type:             r.Form.Get("type"),
		Question:         r.Form.Get("question"),
		DifficultyLevel:  difficultyLevel,
		TestID:           testID,
		UserID:           userID,
		TopicSelectionID: topicSelectionID,
	}

	// 2) convert answerChoices model fields
	// need to iterate through each answer_choice - user can add as many as they want or remove them
	// for each answer_choice we need to make a new post request

	// we do not need a separate form for answer_choices because it is all in the same form, we are just separating the fields on the backend (here)
	form := forms.New(r.PostForm)

	form.Required("type", "question", "difficulty_level", "test_id", "user_id", "topic_selection_id")
	// we will get the question_id from a return value from the database
	// we need to validate we have everything we need before we insert the question into the database (both question and answer_choices)
	// we will then insert the answer_choices into the database if the form is valid and the question is inserted successfully
	form.Required("choice", "correct")

	// this is so we can repopulate the form on errors, so they do not have to redo the form
	// if !form.Valid() {
	// 	data := make(map[string]interface{})
	// 	data["question"] = question
	// 	data["answerChoices"] = answerChoices

	// 	render.Template(w, r, "make-question.page.tmpl", &models.TemplateData{
	// 		Form: form,
	// 		Data: data,
	// 	})
	// 	return
	// }

	// we need to insert our question into the database
	err = m.DB.InsertQuestion(question)
	if err != nil {
		log.Println(err)
		return
	}

	// m.App.Session.Put(r.Context(), "test", test) // don't need this

	// we want to redirect so our user does not submit twice (this is crucial for financial posts)
	http.Redirect(w, r, "/make-question", http.StatusSeeOther) // we can use other status' for redirect, but this one works well enough (303)
}

// Register displays the registration page
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "register.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostRegister handles the registration form
func (m *Repository) PostRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Form submission error")
		http.Redirect(w, r, "/user/register", http.StatusTemporaryRedirect)
		return
	}

	// create a new user
	user := models.User{
		FirstName:   r.Form.Get("first_name"),
		LastName:    r.Form.Get("last_name"),
		Email:       r.Form.Get("email"),
		Password:    r.Form.Get("password"),
		AccessLevel: 1,
	}
	fmt.Println("here")

	// validate the form
	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email", "password")
	form.IsEmail("email")
	form.MinLength("password", 8, r)
	if !form.Valid() {
		render.Template(w, r, "register.page.tmpl", &models.TemplateData{Form: form})
		return
	}

	fmt.Println("here 1")

	// insert the user into the database
	err = m.DB.RegisterUser(user)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Server error. Please try again.")
		http.Redirect(w, r, "/user/register", http.StatusSeeOther)
		return
	}

	// send a flash success message to the user and redirect to login page
	m.App.Session.Put(r.Context(), "flash", "Account created. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context()) // stops a session fixation attack by renewing the token every time the user logs in
	// should be done every time a user logs in or logs out

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{Form: form})
		return
	}
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	user := models.User{
		ID: id,
	}

	// this session is what we are going to use to check if the user is logged in or not
	m.App.Session.Put(r.Context(), "user", user)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())    // destroys the session
	_ = m.App.Session.RenewToken(r.Context()) // good practice to renew the token after destroying the session
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) ShowAccountDashboard(w http.ResponseWriter, r *http.Request) {
	u, ok := m.App.Session.Get(r.Context(), "user").(models.User)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Could not get user from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// get the user from the database
	user, err := m.DB.GetUserByID(u.ID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Could not get user from database")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["user"] = user

	render.Template(w, r, "account-dashboard.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ShowAccountEdit(w http.ResponseWriter, r *http.Request) {
	u, ok := m.App.Session.Get(r.Context(), "user").(models.User)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Could not get user from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user, err := m.DB.GetUserByID(u.ID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Could not get user from database")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["user"] = user

	render.Template(w, r, "account-edit.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) PostAccountEdit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Could not parse form")
		http.Redirect(w, r, "/account/edit", http.StatusSeeOther)
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("user_id"))

	user := models.User{
		ID:        userID,
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = user

		render.Template(w, r, "account-edit.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.UpdateUser(user)
	if err != nil {
		log.Println(err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, "/account/dashboard", http.StatusSeeOther)
}

// JSON response structure for account delete
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AccountDelete deletes a user account
func (m *Repository) PostAccountDelete(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("user_id"))
	err = m.DB.DeleteUser(userID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error deleting the user",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonResponse{
		OK:      true,
		Message: "User successfully deleted",
	}
	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
