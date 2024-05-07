package dbrepo

import (
	"database/sql"

	"github.com/Rich-Wilkyness/kether/internal/config"
	"github.com/Rich-Wilkyness/kether/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates a new postgres repository
// if we decide to use a different database, we can create a new func to handle it
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
