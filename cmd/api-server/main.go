package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"api-server/internal/config"
	"api-server/internal/database"
	tracing "api-server/internal/observability"
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

	// Initialize OpenTelemetry
	log.Printf("Initializing OpenTelemetry with service name '%s' and endpoint '%s'",
		cfg.ServiceName, cfg.OtlpEndpoint)
	shutdown := tracing.InitTracer(cfg.ServiceName, cfg.OtlpEndpoint)
	defer func() {
		log.Println("Shutting down OpenTelemetry")
		shutdown()
	}()

	// Send Initial span on startup
	go func() {
		time.Sleep(2 * time.Second) // Wait for server to start
		tracer := otel.Tracer("startup-inital-span")
		ctx := context.Background()
		_, span := tracer.Start(ctx, "startup-initial-span")
		log.Println("Created initial span at startup")
		time.Sleep(100 * time.Millisecond)
		span.End()
		log.Println("Ended initial span after startup")
	}()

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

	// Wrap the router with OpenTelemetry middleware
	handler := otelhttp.NewHandler(r, "api-server")

	// Start server
	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
