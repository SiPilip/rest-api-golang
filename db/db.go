package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		log.Fatal("DB_DSN is not set in .env file")
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Hour)

	createTables()

	err = DB.Ping()
	if err != nil {
		panic("Error pinging database:" + err.Error())
	}

	fmt.Println("Database connected successfully")
}

func createTables() {
	createusersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY AUTO_INCREMENT,
		email VARCHAR(255) NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`

	_, err := DB.Exec(createusersTable)
	if err != nil {
		panic("Error creating users table:" + err.Error())
	}
	
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id INT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createEventsTable)
	if err != nil {
		panic("Error creating events table:" + err.Error())
	}

	createRegistrationTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id INT PRIMARY KEY AUTO_INCREMENT,
		event_id INT NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(event_id) REFERENCES events(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createRegistrationTable)
	if err != nil {
		panic("Error creating registrations table:" + err.Error())
	}
}