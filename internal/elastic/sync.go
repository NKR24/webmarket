package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"shop/internal/product"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSync struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticSync(elasticHost string) (*ElasticSync, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticSync{
		client: client,
		index:  "products",
	}, nil
}

func (e *ElasticSync) IndexProduct(product *product.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("error marshaling product to JSON: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      e.index,
		DocumentID: strconv.Itoa(product.ID),
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return fmt.Errorf("error indexing product in Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error indexing product ID %d: %s", product.ID, res.String())
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	log.Printf("Product ID %d indexed successfully", product.ID)
	return nil
}

func (e *ElasticSync) DeleteProduct(productID int) error {
	req := esapi.DeleteRequest{
		Index:      e.index,
		DocumentID: strconv.Itoa(productID),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return fmt.Errorf("error deleting product from Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error deleting product ID %d: %s", productID, res.String())
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	log.Printf("Product ID %d deleted successfully", productID)
	return nil
}
