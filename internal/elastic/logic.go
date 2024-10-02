package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func (es *ElsaticClient) IndexProduct(index string, docID string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling product data: %v", err)
		return err
	}

	log.Printf("Indexing product in Elasticsearch. Index: %s, DocID: %s, Data: %s", index, docID, jsonData)

	res, err := es.es.Index(index, bytes.NewReader(jsonData), es.es.Index.WithDocumentID(docID))
	if err != nil {
		log.Printf("Error sending request to Elasticsearch: %v", err)
		return err
	}
	defer res.Body.Close()

	log.Printf("Elasticsearch response status: %s", res.Status())

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Printf("Elasticsearch error response: %s", string(bodyBytes))
		return fmt.Errorf("Elasticsearch indexing error: %s", string(bodyBytes))
	}

	log.Println("Product successfully indexed in Elasticsearch")
	return nil
}

func (es *ElsaticClient) DeleteProduct(index, id string) error {
	res, err := es.es.Delete(index, id)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil

}
