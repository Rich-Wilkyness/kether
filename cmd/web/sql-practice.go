package main

import (
	"database/sql"
	"fmt"
	"log"

	// this must be imported like this in order to connect to the db
	_ "github.com/jackc/pgx/v5/stdlib"
)

func sqlPractice() {
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=postgres user=postgres password=lebrum1203")
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect: %v\n", err))
	}
	defer conn.Close()

	log.Println("Connected to db!")

	// test connection
	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot ping db!")
	}
	log.Println("Pinged db!")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// insert row. use `` backticks so you can format your query over multiple lines.
	query := `
		INSERT into users (first_name, last_name) values ($1, $2)
	`
	_, err = conn.Exec(query, "Jack", "Brown") // we use $value as substitutes becuase it is a safer way to make a query
	// forces db to not allow arbitrary strings (sql injection) to be passed in
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted row!")
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// update row
	stmt := `
		UPDATE users SET first_name = $1 WHERE first_name = $2
	`
	_, err = conn.Exec(stmt, "John", "Jack")
	if err != nil {
		log.Fatal(err)
	}

	// get 1 row by id
	query = `select id, first_name, last_name from users where id = $1`
	var firstName, lastName string
	var id int

	row := conn.QueryRow(query, 1)
	err = row.Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Row is", id, firstName, lastName)

	// delete row
	query = `DELETE from users where id = $1`
	_, err = conn.Exec(query, 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Deleted a row!")
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("select id, first_name, last_name from users")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close() // remember to close, takes resources and security risk

	var firstName, lastName string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("Record is", id, firstName, lastName)
	}
	if err = rows.Err(); err != nil { // this is just good practice to make sure you catch an error if something happened while scanning
		log.Fatal("Error scanning rows", err)
	}
	fmt.Println("-----------------------------------------")
	return nil
}

// ---------------------------- setting up connection to db --------------------------
// 1. go https://github.com/jackc/pgx/wiki/Getting-started-with-pgx-through-database-sql
// 2. follow instructions
// 3. make sure you have the correct imports
