package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rich-Wilkyness/kether/internal/models"
	"golang.org/x/crypto/bcrypt"
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

// RegisterUser registers a user in the database
func (m *postgresDBRepo) RegisterUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	fmt.Println("hashed password:")

	stmt := `INSERT into users 
				(first_name, last_name, email, password, access_level, created_at, updated_at)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7)`

	_, err = m.DB.ExecContext(ctx, stmt, u.FirstName, u.LastName, u.Email, string(hashedPassword), u.AccessLevel, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

// returns user id
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, first_name, last_name, email, created_at, updated_at
		FROM users 
		WHERE id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, email = $3, updated_at = $4
	`
	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user from the database
func (m *postgresDBRepo) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE 
			FROM users 
			WHERE id = $1`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	query := `
		SELECT id, password 
		FROM users 
		where email = $1
	`
	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
