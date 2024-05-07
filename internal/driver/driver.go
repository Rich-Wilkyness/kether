package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"        // pgx driver
	_ "github.com/jackc/pgx/v5"        // pgx driver
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
)

// this file is where we connect to the database and create a new connection pool
// we also define the methods that will be used to interact with the database

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB // reason for using a struct: we might have multiple different types of databases in the future
	// by making it a struct we can add more fields to the struct to support other databases
	// we will also be adding things to the repository to support other databases
}

var dbConn = &DB{} // brackets is because we are creating a new instance of the struct, which will be empty

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifetime = 5 * time.Minute

// create connection pool
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	// these are the settings for the connection pool
	// they keep the db from being overwhelmed with too many connections
	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetMaxIdleConns(maxIdleDBConn)
	d.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// test database connection
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil { // remember ; means and in go
		return nil, err
	}

	return db, nil
}
