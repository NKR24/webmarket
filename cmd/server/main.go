package main

import (
	"fmt"
	"log"
	"os"

	"market/internal/db"
	"market/internal/elastic"
	"market/internal/kafka"
	"market/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем параметры подключения к PostgreSQL из переменных окружения
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	// Формируем строку подключения к базе данных
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Подключаемся к PostgreSQL
	db, err := db.Connect(dsn)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}

	// Подключаем Kafka Producer
	kp := kafka.NewProducer(os.Getenv("KAFKA_BROKER"), "products")

	p := new(models.Product)
	// Создаем продукт и отправляем событие в Kafka
	p.Name = "Product 1"
	p.Description = "A test product"
	p.Price = 19.99
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
