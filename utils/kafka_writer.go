package utils

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter wraps the Kafka writer
type KafkaWriter struct {
	Writer *kafka.Writer
}

// InitializeKafkaWriter initializes a Kafka writer for the given topic
func InitializeKafkaWriter(brokerAddress, topic string) *KafkaWriter {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	log.Printf("Initialized Kafka writer for topic: %s\n", topic)
	return &KafkaWriter{Writer: writer}
}

// WriteMessage writes a message to Kafka
func (kw *KafkaWriter) WriteMessage(key, value string) error {
	return kw.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		})
}

// Close closes the Kafka writer
func (kw *KafkaWriter) Close() {
	if kw.Writer != nil {
		kw.Writer.Close()
	}
}
