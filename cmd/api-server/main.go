package main

import (
	"log"
	"net/http"

	"api-server/internal/config"
	"api-server/internal/database"
	"api-server/internal/routes"
	"api-server/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Kafka producer with authentication
	services.InitKafkaProducer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		cfg.KafkaUsername,
		cfg.KafkaPassword,
		cfg.KafkaAuth,
	)
	defer services.CloseKafkaProducer()

	// Register routes
	r := routes.RegisterRoutes()

	// Start server
	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
