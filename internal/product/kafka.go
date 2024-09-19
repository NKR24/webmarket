package product

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	prducer *kafka.Producer
	topic   string
}

func NewKafkaProducer(brokers, topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}
	defer p.Close()
	return &KafkaProducer{prducer: p, topic: topic}, err
}

func (kp *KafkaProducer) SendProductMessage(product *Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	err = kp.prducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kp.topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}
