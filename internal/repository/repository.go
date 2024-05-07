package repository

import (
	"github.com/Rich-Wilkyness/kether/internal/models"
)

// Purpose of this file is to give access to the database to the handlers
type DatabaseRepo interface {
	AllUsers() bool

	InsertTest(res models.Test) error
	InsertQuestion(r models.Question) error
}
