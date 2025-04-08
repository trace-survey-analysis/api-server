package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

// Metadata for an uploaded trace survey
type TraceUploadMessage struct {
	TraceID      string    `json:"traceId"`
	CourseID     string    `json:"courseId"`
	FileName     string    `json:"fileName"`
	GCSBucket    string    `json:"gcsBucket"`
	GCSPath      string    `json:"gcsPath"`
	InstructorID string    `json:"instructorId"`
	SemesterTerm string    `json:"semesterTerm"`
	Section      string    `json:"section"`
	UploadedBy   string    `json:"uploadedBy"`
	UploadedAt   time.Time `json:"uploadedAt"`
}

// New Kafka producer
func NewProducer(brokers []string, topic, username, password string, enableAuth bool) (*Producer, error) {
	// Basic writer configuration
	writerConfig := kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		RequiredAcks: int(kafka.RequireAll),
		WriteTimeout: 10 * time.Second,
		BatchSize:    100,
		BatchTimeout: 1 * time.Millisecond,
		// Logger:       kafka.LoggerFunc(logf), // Uncomment for debugging
	}

	// Add authentication if enabled
	if enableAuth {
		if username == "" || password == "" {
			return nil, fmt.Errorf("kafka authentication enabled but missing credentials")
		}

		mechanism := plain.Mechanism{
			Username: username,
			Password: password,
		}

		dialer := &kafka.Dialer{
			Timeout:       10 * time.Second,
			DualStack:     true,
			SASLMechanism: mechanism,
			// TLS:           &tls.Config{},
		}

		writerConfig.Dialer = dialer
		log.Println("Kafka SASL authentication enabled")
	}

	writer := kafka.NewWriter(writerConfig)

	return &Producer{
		writer: writer,
		topic:  topic,
	}, nil
}

// Send a trace survey upload notification to Kafka
func (p *Producer) PublishTraceUpload(ctx context.Context, message TraceUploadMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(message.TraceID),
		Value: data,
		Time:  time.Now(),

		Headers: []kafka.Header{
			{Key: "content-type", Value: []byte("application/json")},
			{Key: "source", Value: []byte("api-server")},
		},
	})

	if err != nil {
		return fmt.Errorf("error writing message to Kafka: %w", err)
	}

	log.Printf("Published trace upload message to Kafka for trace ID: %s", message.TraceID)
	return nil
}

// Close the Kafka writer
func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("error closing Kafka writer: %w", err)
	}
	return nil
}

// Helper function for debugging
func logf(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
