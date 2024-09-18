package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewProducer(broker, topic string) *KafkaProducer {
	return &KafkaProducer{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{broker},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (kp *KafkaProducer) SendMessage(event string, value any) error {
	message, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = kp.Writer.WriteMessages(context.TODO(), kafka.Message{
		Key:   []byte(event),
		Value: message,
	})

	if err != nil {
		log.Println("Error writing message: ", err)
	}

	return err
}
