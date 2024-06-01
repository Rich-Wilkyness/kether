# Overview

My goal with this program is to practice implementing my learning of Golang, Bootstrap, Databases, and SQL. My secondary goal is to develop a web app for parents who want to develop their own home-school curriculum. 

End goal is for this program to allow one to create a test, add questions to that test, and retrieve a test that will populate with the assigned questions. 
Currently, the program allows for someone to register as a user, login/out, and access their personal account to edit or delete it with authentication. 

# How To Run

1. Download the repository
2. Ensure all dependencies are downloaded - you can do this by running go mod tidy and then run 'go get (dependency)' in the terminal if certain dependencies were not resolved - ex: go get github.com/Rich-Wilkyness/kether
2. Go to root directory and run "go run ./cmd/web" in the terminal

[Software Demo Video](https://youtu.be/VSK3z5hsX68)

# Development Environment

1. [Golang](https://go.dev/) - backend
2. [Boostrap](https://getbootstrap.com/docs/5.3/getting-started/introduction/) - frontend 
3. [PostgreSQL](https://www.postgresql.org/download/) - for the database
4. [DBeaver](https://dbeaver.io/download/) - GUI for database
5. https://gobuffalo.io/documentation/database/fizz/ - database migration
6. and many other dependencies in the go.mod file for authentication, routing, etc. 

# Useful Websites

* [Udemy](https://www.udemy.com/user/trevor-sawler/)
Trevor Sawler has several top tier Golang courses.