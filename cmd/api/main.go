package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"shop/config"
	"shop/internal/product"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := sql.Open("postgres",
		"host="+cfg.DBHost+" port="+cfg.DBPort+" user="+cfg.DBUser+" password="+cfg.DBPassword+" dbname="+cfg.DBName+" sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	kafkaProducer, err := product.NewKafkaProducer(cfg.KafkaBrokers, "products")
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	handler := product.NewHandler(db, kafkaProducer)

	r := mux.NewRouter()
	r.HandleFunc("/products", handler.CreateProduct).Methods("POST")
	r.HandleFunc("/products", handler.GetAllProducts).Methods("GET")
	r.HandleFunc("/products/{id}", handler.GetProductByID).Methods("GET")
	r.HandleFunc("/products/{id}", handler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handler.DeleteProduct).Methods("DELETE")

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "9090"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
