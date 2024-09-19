package main

import (
	"encoding/json"
	"log"
	"os"

	"shop/internal/elastic"
	"shop/internal/product"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKERS"),
		"group.id":          "shop-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatal(err)
	}

	consumer.SubscribeTopics([]string{"products"}, nil)

	elasticSync, err := elastic.NewElasticSync(os.Getenv("ELASTIC_HOST"))
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %v", err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		var product product.Product
		if err := json.Unmarshal(msg.Value, &product); err != nil {
			log.Println("Failed to unmarshal product:", err)
			continue
		}

		// Вставляем или обновляем продукт в Elasticsearch
		if err := elasticSync.IndexProduct(&product); err != nil {
			log.Printf("Error indexing product with ID %d: %v", product.ID, err)
		} else {
			log.Printf("Product with ID %d successfully synced to Elasticsearch", product.ID)
		}
	}
}
