package elastic

import (
	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticClient struct {
	Client *elasticsearch.Client
}

func NewElasticClient(address string) (*ElasticClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{address},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ElasticClient{Client: es}, nil
}

func (es *ElasticClient) IndexProduct(index string, docID string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(jsonData)

	res, err := es.Client.Index(index, reader, es.Client.Index.WithDocumentID(docID))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
