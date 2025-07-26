// kafka_service.go
//
// KafkaService handles producing and consuming chat messages via Kafka.
//
// Workflow:
// - When a new chat message is received via REST API, it is published to the Kafka topic "chat-messages".
// - A Kafka consumer subscribes to this topic, and for each new message, broadcasts it to connected WebSocket clients based on session/room/user ID.
// - This decouples message delivery, improves scalability, and enables future extensibility (e.g., analytics, bots).

package service

import (
	"context"
	"encoding/json"
	"os"

	"github.com/segmentio/kafka-go"
)

// KafkaService manages Kafka producer and consumer for chat messages.
type KafkaService struct {
	Writer *kafka.Writer
	// Reader *kafka.Reader
	Topic string
}

// NewKafkaService initializes a new KafkaService.
func NewKafkaService() *KafkaService {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "chat-messages"
	}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	// reader := kafka.NewReader(kafka.ReaderConfig{
	// 	Brokers:   []string{broker},
	// 	Topic:     topic,
	// 	GroupID:   "livechat-ws-group",
	// 	Partition: 0,
	// 	MinBytes:  1,
	// 	MaxBytes:  10e6,
	// })
	return &KafkaService{
		Writer: writer,
		// Reader: reader,
		Topic: topic,
	}
}

// PublishMessage publishes a chat message to Kafka.
func (k *KafkaService) PublishMessage(ctx context.Context, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return k.Writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}

// ConsumeAndBroadcast consumes messages from Kafka and broadcasts to WebSocket clients.
// func (k *KafkaService) ConsumeAndBroadcast(ctx context.Context, broadcastFunc func([]byte)) {
// 	for {
// 		m, err := k.Reader.ReadMessage(ctx)
// 		if err != nil {
// 			log.Printf("Kafka consume error: %v", err)
// 			time.Sleep(time.Second)
// 			continue
// 		}
// 		broadcastFunc(m.Value)
// 	}
// }

// Close closes Kafka connections.
func (k *KafkaService) Close() error {
	// if err := k.Writer.Close(); err != nil {
	// 	return err
	// }
	// return k.Reader.Close()
	return k.Writer.Close()
}
