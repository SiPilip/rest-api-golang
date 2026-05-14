package main

import (
	"REST-API/db"
	"REST-API/routes"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Setup slog
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	
	// Load .env
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file:", "error", err)
		os.Exit(1)
	}
	
	// Init database
	db.InitDB()

	// Setup gin
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Static("/uploads", "./uploads")

	// GRACEFUL SHUTDOWN
	port := os.Getenv("SERVER_PORT")
	srv := &http.Server{
		Addr: ":" + port,
		Handler: server,
	}

	// Server start
	go func() {
		slog.Info("Server started", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Tunggu sinyal shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Beri waktu 5 detik nutuk selesaikan request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	// Tutup koneksi database
	if err := db.DB.Close(); err != nil {
		slog.Error("Error closing database", "error", err)
	}	

	slog.Info("Server gracefully stopped")
}