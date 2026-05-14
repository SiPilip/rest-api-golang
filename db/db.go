package db

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		slog.Error("DB_DSN is not set in .env file")
		os.Exit(1)
	}

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		slog.Error("Error connecting to database:", "error", err.Error())
		os.Exit(1)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Hour)

	err = DB.Ping()
	if err != nil {
		panic("Error pinging database:" + err.Error())
	}

	slog.Info("Database connected successfully")
}