package services

import (
	"api-server/internal/kafka"
	"log"
	"sync"
)

var (
	kafkaProducer     *kafka.Producer
	kafkaProducerLock sync.RWMutex
)

// Initialize the Kafka producer
func InitKafkaProducer(brokers []string, topic, username, password string, enableAuth bool) {
	kafkaProducerLock.Lock()
	defer kafkaProducerLock.Unlock()

	if len(brokers) == 0 {
		log.Println("No Kafka brokers configured, skipping producer initialization")
		return
	}

	var err error
	kafkaProducer, err = kafka.NewProducer(brokers, topic, username, password, enableAuth)
	if err != nil {
		log.Printf("Failed to initialize Kafka producer: %v", err)
		return
	}

	log.Printf("Kafka producer initialized successfully for topic: %s", topic)
}

// Return the initialized Kafka producer
func GetKafkaProducer() *kafka.Producer {
	kafkaProducerLock.RLock()
	defer kafkaProducerLock.RUnlock()
	return kafkaProducer
}

// Close the Kafka producer
func CloseKafkaProducer() {
	kafkaProducerLock.Lock()
	defer kafkaProducerLock.Unlock()

	if kafkaProducer != nil {
		err := kafkaProducer.Close()
		if err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		} else {
			log.Println("Kafka producer closed successfully")
		}
		kafkaProducer = nil
	}
}
