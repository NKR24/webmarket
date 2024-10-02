package consumer

import (
	"context"
	"encoding/json"
	"log"
	"shop/internal/elastic"
	"shop/internal/product"
)

func ConsumeMessage(broker, topic, index, address string, parition, minbytes, maxbytes int) {
	reader, err := NewKafkaReader(broker, topic, parition, minbytes, maxbytes)
	if err != nil {
		log.Printf("Failed to create new kafka reader: %v", err)
	}

	es, err := elastic.NewElasticClient(address)
	if err != nil {
		log.Printf("Failed to create elastic client: %v", err)
	}

	for {
		msg, err := reader.kr.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message from Kafka:", err)
			continue
		}

		log.Printf("Message received from Kafka: %s", string(msg.Value))

		p := new(product.Product)
		err = json.Unmarshal(msg.Value, &p)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		err = es.IndexProduct(index, string(msg.Key), p)
		if err != nil {
			log.Println("Error indexing product in Elasticsearch:", err)
		} else {
			log.Println("Product successfully indexed in Elasticsearch")
		}
	}
}
