package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB()  {
	var err error
	DB, err = sql.Open("sqlite3", "event-booking.db")

	if err != nil {
		basic_err_msg := "Error: Could not connect to database!"
		fmt.Printf("%v %v\n", basic_err_msg, err)
		panic(basic_err_msg)
	}

	fmt.Println("Successfully opened database!")

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTable()
}

func createTable()  {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		password TEXT
	)
	`

	_, err := DB.Exec(createUserTable)

	if err != nil {
		basic_err_msg := "Error: Could not create 'users' table!"
		fmt.Printf("%v %v\n", basic_err_msg, err)
		panic(basic_err_msg)
	}

	createEventTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createEventTable)

	if err != nil {
		basic_err_msg := "Error: Could not create 'events' table!"
		fmt.Printf("%v %v\n", basic_err_msg, err)
		panic(basic_err_msg)
	}

	createRegistrationTable := `
	CREATE TABLE IF NOT EXISTS registrations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_id INTEGER UNIQUE,
	user_id INTEGER,
	FOREIGN KEY(event_id) REFERENCES events(id),
	FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createRegistrationTable)

	if err != nil {
		basic_err_msg := "Error: Could not create 'registrations' table!"
		fmt.Printf("%v %v\n", basic_err_msg, err)
		panic(basic_err_msg)
	}
	
	fmt.Println("Successfully created tables!")
}