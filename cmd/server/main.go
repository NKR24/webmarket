package main

import (
	"log"
	"market/internal/db"
	"market/internal/elastic"
	"market/internal/kafka"
	"market/internal/models"
	"os"
)

func main() {
	// Подключаем PostgreSQL
	db, err := db.Connect(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	// Подключаем Kafka Producer
	kp := kafka.NewProducer(os.Getenv("KAFKA_BROKER"), "products")

	// Создаем продукт и отправляем событие в Kafka
	p := &models.Product{Name: "Product 1", Description: "A test product", Price: 19.99}
	_, err = db.CreateProductWithKafka(p, kp)
	if err != nil {
		log.Fatal("Failed to create product:", err)
	}

	// Подключаем Elasticsearch
	es, err := elastic.NewElasticClient(os.Getenv("ELASTICSEARCH_URL"))
	if err != nil {
		log.Fatal("Failed to connect to Elasticsearch:", err)
	}

	// Консьюмер для Kafka
	go kafka.ConsumeMessages(os.Getenv("KAFKA_BROKER"), "products", *es, "products")

	select {}
}
