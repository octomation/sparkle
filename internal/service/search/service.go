package search

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"go.octolab.org/pointer"

	"go.octolab.org/ecosystem/sparkle/internal/service/search/schema"
)

func New(client *typesense.Client) *Service {
	return &Service{client: client}
}

type Service struct {
	client *typesense.Client
}

func (service *Service) Find(query string) ([]schema.Sparkle, error) {
	ctx := context.TODO()

	params := &api.SearchCollectionParams{
		Q:                   query,
		QueryBy:             schema.Sparkle{}.QueryFields(),
		TypoTokensThreshold: pointer.ToInt(5),
		UseCache:            pointer.ToBool(false),
	}

	result, err := service.client.
		Collection(schema.Sparkle{}.Collection()).
		Documents().
		Search(ctx, params)
	if err != nil {
		return nil, err
	}
	if result.Hits == nil || len(*result.Hits) == 0 {
		return nil, errNotFound
	}

	buf := bytes.NewBuffer(nil)
	docs := make([]schema.Sparkle, 0, len(*result.Hits))
	for _, hit := range *result.Hits {
		buf.Reset()
		if err := json.NewEncoder(buf).Encode(hit.Document); err != nil {
			return nil, err
		}

		var doc schema.Sparkle
		if err := json.NewDecoder(buf).Decode(&doc); err != nil {
			return nil, err
		}
		doc.Highlights = *hit.Highlights
		docs = append(docs, doc)
	}
	return docs, nil
}

func (service *Service) Index(docs ...schema.Sparkle) error {
	ctx := context.TODO()

	collection := service.client.
		Collection(schema.Sparkle{}.Collection()).
		Documents()
	for _, doc := range docs {
		if _, err := collection.Upsert(ctx, doc); err != nil {
			return err
		}
	}
	return nil
}
