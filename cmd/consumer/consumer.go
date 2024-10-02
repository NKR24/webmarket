package consumer

import (
	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	kr *kafka.Reader
}

func NewKafkaReader(broker, topic string, partition, minbytes, maxbytes int) (*KafkaReader, error) {
	config := kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		Partition: partition,
		MinBytes:  minbytes,
		MaxBytes:  maxbytes,
	}

	reader := kafka.NewReader(config)

	return &KafkaReader{kr: reader}, nil
}

func (r *KafkaReader) Close() error {
	return r.kr.Close()
}
