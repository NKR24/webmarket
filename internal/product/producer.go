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

func NewKafkaWriter(broker, topic string) (*KafkaWriter, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaWriter{kw: writer}, nil
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

func (w *KafkaWriter) CreateKafkaTopic(broker, topic string, partition int) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, partition)
	if err != nil {
		log.Fatalf("Ошибка в создании KafkaTopic: %v", err)
	}
	defer conn.Close()

	return nil
}
