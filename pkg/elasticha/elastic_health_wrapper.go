package elasticha

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

type ClientWrapper struct {
	client *elasticsearch.TypedClient
}

func NewClientWrapper(config elasticsearch.Config) (*ClientWrapper, error) {
	client, err := elasticsearch.NewTypedClient(config)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
		return nil, err
	}

	return &ClientWrapper{client: client}, nil
}

// Function to connect to Elasticsearch with default configuration
func ConnectElasticsearch() (*ClientWrapper, error) {
	// Define the configuration
	config := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
			//
		},
		//
	}

	// Create a new client wrapper
	clientWrapper, err := NewClientWrapper(config)
	if err != nil {
		log.Fatalf("Error creating client wrapper: %v", err)
		return nil, err
	}

	return clientWrapper, nil
}
