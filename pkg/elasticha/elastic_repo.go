package elasticha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
)

type ServerDocument struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Ipv4   string    `json:"ipv4"`
	Status string    `json:"status"`
	At     time.Time `json:"at"`
	In     int       `json:"in"`
}

type Repository struct {
	client *ClientWrapper
	index  string
}

func NewRepository(client *ClientWrapper, index string) *Repository {
	return &Repository{
		client: client,
		index:  index,
	}
}

func (r *Repository) CreateOne(item ServerDocument) (ServerDocument, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return item, err
	}

	req := esapi.IndexRequest{
		Index:      r.index,
		DocumentID: item.ID.String(),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), r.client.client)
	if err != nil {
		return item, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return item, fmt.Errorf("error indexing document ID=%s", item.ID.String())
	}

	return item, nil
}

func (r *Repository) FindOneById(id uuid.UUID) (ServerDocument, error) {
	var item ServerDocument
	req := esapi.GetRequest{
		Index:      r.index,
		DocumentID: id.String(),
	}

	res, err := req.Do(context.Background(), r.client.client)
	if err != nil {
		return item, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return item, fmt.Errorf("error getting document ID=%s", id.String())
	}

	if err := json.NewDecoder(res.Body).Decode(&item); err != nil {
		return item, err
	}

	return item, nil
}

func (r *Repository) UpdateOneById(id uuid.UUID, items map[string]interface{}) (ServerDocument, error) {
	var updatedItem ServerDocument
	data, err := json.Marshal(map[string]interface{}{
		"doc": items,
	})
	if err != nil {
		return updatedItem, err
	}

	req := esapi.UpdateRequest{
		Index:      r.index,
		DocumentID: id.String(),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), r.client.client)
	if err != nil {
		return updatedItem, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return updatedItem, fmt.Errorf("error updating document ID=%s", id.String())
	}

	return r.FindOneById(id)
}

func (r *Repository) DeleteOneById(id uuid.UUID) error {
	req := esapi.DeleteRequest{
		Index:      r.index,
		DocumentID: id.String(),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), r.client.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting document ID=%s", id.String())
	}

	return nil
}

func (r *Repository) FindOneByAttribute(key string, value interface{}) (ServerDocument, error) {
	var item ServerDocument
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				key: value,
			},
		},
	}

	data, err := json.Marshal(query)
	if err != nil {
		return item, err
	}

	req := esapi.SearchRequest{
		Index: []string{r.index},
		Body:  bytes.NewReader(data),
	}

	res, err := req.Do(context.Background(), r.client.client)
	if err != nil {
		return item, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return item, errors.New("error searching document")
	}

	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source ServerDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return item, err
	}

	if len(searchResult.Hits.Hits) > 0 {
		item = searchResult.Hits.Hits[0].Source
	} else {
		return item, errors.New("no document found")
	}

	return item, nil
}
