package repository

import (
	"github.com/Rich-Wilkyness/kether/internal/models"
)

// Purpose of this file is to give access to the database to the handlers
type DatabaseRepo interface {
	AllUsers() bool

	// Test methods
	InsertTest(res models.Test) error
	InsertQuestion(r models.Question) error

	// User methods
	RegisterUser(u models.User) error
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	DeleteUser(id int) error
	Authenticate(email, testPassword string) (int, string, error)
}
