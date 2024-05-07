package handlers

import (
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
	if !form.Valid() {
		data := make(map[string]interface{})
		data["question"] = question
		data["answerChoices"] = answerChoices

		render.Template(w, r, "make-question.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

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
