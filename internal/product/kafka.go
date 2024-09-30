package product

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	kw *kafka.Writer
}

// func CreateKafkaTopic(broker, topic string, partition int) error {
// 	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, partition)
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()

// 	// _, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	// defer cancel()

// 	// topicConfig := kafka.TopicConfig{
// 	// 	Topic:             topic,
// 	// 	NumPartitions:     1,
// 	// 	ReplicationFactor: 1,
// 	// }

// 	// err = conn.CreateTopics(topicConfig)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// log.Printf("Kafka topic '%s' created successfully", topic)
// 	return nil
// }

func NewKafkaWriter(broker, topic string, partition int) *KafkaWriter {
	// err := CreateKafkaTopic(broker, topic, partition)
	// if err != nil {
	// 	log.Printf("Failed to create Kafka topic: %v", err)
	// 	return nil, err
	// }

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	return &KafkaWriter{kw: writer}
}

func (w *KafkaWriter) SendMessage(p *Product) error {
	message, err := json.Marshal(p)
	if err != nil {
		return err
	}

	err = w.kw.WriteMessages(context.Background(), kafka.Message{
		Value: message,
		Time:  time.Now(),
	})
	if err != nil {
		log.Printf("Error sending message to Kafka: %v", err)
		return err
	}
	log.Printf("Product %v sent to Kafka", p)
	return nil
}
