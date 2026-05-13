package main

import (
	"REST-API/db"
	"REST-API/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":"+os.Getenv("SERVER_PORT"))
}