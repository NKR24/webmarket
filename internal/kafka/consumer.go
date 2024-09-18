package kafka

import (
	"context"
	"encoding/json"
	"log"
	"market/internal/elastic"
	"market/internal/models"

	"github.com/segmentio/kafka-go"
)

func ConsumeMessages(broker, topic string, esClient elastic.ElasticClient, index string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "product-group",
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message from Kafka:", err)
			continue
		}

		p := new(models.Product)
		err = json.Unmarshal(msg.Value, &p)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		err = esClient.IndexProduct(index, string(msg.Key), p)
		if err != nil {
			log.Println("Error indexing product in Elasticsearch:", err)
		}
	}
}
