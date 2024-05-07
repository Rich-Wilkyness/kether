package dbrepo

import (
	"context"
	"time"

	"github.com/Rich-Wilkyness/kether/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertTest(res models.Test) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // this solves the problem of hanging connections (user losing connection to internet, etc.)
	defer cancel()

	// we use stmt when we are putting data into the database
	// we use query when we are getting data from the database
	stmt := `INSERT into tests (name, version, class_id, user_id, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)`

	// the .Exec method has risk. User loses connection to internet, closes app, etc. We need to handle this
	// we use context to handle this (ExecContext or instead QueryRowContext in this case because we are returning a value, not just inserting data)
	_, err := m.DB.ExecContext(ctx, stmt, res.Name, res.Version, res.ClassID, res.UserID, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) InsertQuestion(r models.Question) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT into questions
				(type, question, difficulty_level, test_id, user_id, topic_selection_id, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := m.DB.ExecContext(ctx, stmt, r.Type, r.Question, r.DifficultyLevel, r.TestID, r.UserID, r.TopicSelectionID, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}
