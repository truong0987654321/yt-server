package main

import (
	"log"
	"yt-clone-server/internal/config"
	"yt-clone-server/internal/repository/postgres"
	"yt-clone-server/internal/router"
)

func main() {
	log.Println("Starting the Go Backend Server...")

	// Load config
	cfg := config.Load()

	// Kết nối DB
	db, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Database connection established successfully")

	// Setup router
	r := router.Setup(router.Dependencies{})
	log.Printf("Server is running on port %s (env: %s)", cfg.Port, cfg.Env)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	_ = db
}
