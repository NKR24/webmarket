package elastic

import (
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

type ElsaticClient struct {
	es *elasticsearch.Client
}

func NewElasticClient(address string) (*ElsaticClient, error) {
	config := elasticsearch.Config{
		Addresses: []string{address},
	}
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		log.Printf("Error to start elastic")
		return nil, err
	}

	return &ElsaticClient{es: es}, nil
}
