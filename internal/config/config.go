package config

import (
	"os"
	"strings"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	ServerPort string

	// Kafka configuration
	KafkaBrokers  []string
	KafkaTopic    string
	KafkaUsername string
	KafkaPassword string
	KafkaAuth     bool

	// OpenTelemetry configuration
	ServiceName  string
	OtlpEndpoint string
}

func Load() (*Config, error) {
	// Get Kafka broker list from environment variable
	kafkaBrokersStr := getEnv("KAFKA_BROKERS", "kafka-controller-0.kafka-controller-headless.kafka.svc.cluster.local:9092,kafka-controller-1.kafka-controller-headless.kafka.svc.cluster.local:9092,kafka-controller-2.kafka-controller-headless.kafka.svc.cluster.local:9092")
	kafkaBrokers := strings.Split(kafkaBrokersStr, ",")

	// Get Kafka authentication details
	kafkaUsername := getEnv("KAFKA_USERNAME", "")
	kafkaPassword := getEnv("KAFKA_PASSWORD", "")
	// Enable auth if both username and password are provided
	kafkaAuth := kafkaUsername != "" && kafkaPassword != ""

	return &Config{
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", ""),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// Kafka fields
		KafkaBrokers:  kafkaBrokers,
		KafkaTopic:    getEnv("KAFKA_TOPIC", "trace-survey-uploaded"),
		KafkaUsername: kafkaUsername,
		KafkaPassword: kafkaPassword,
		KafkaAuth:     kafkaAuth,

		// OpenTelemetry fields
		ServiceName:  getEnv("SERVICE_NAME", "api-server"),
		OtlpEndpoint: getEnv("OTLP_ENDPOINT", "localhost:4317"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
