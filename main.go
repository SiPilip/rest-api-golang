package main

import (
	"REST-API/db"
	"REST-API/routes"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file:", "error", err)
		os.Exit(1)
	}
	
	db.InitDB()

	server := gin.Default()
	// Registered static routes
	server.Static("/uploads", "./uploads")

	routes.RegisterRoutes(server)

	server.Run(":"+os.Getenv("SERVER_PORT"))
}