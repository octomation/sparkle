package search

import (
	"context"

	"github.com/typesense/typesense-go/typesense/api"

	"go.octolab.org/ecosystem/sparkle/internal/service/search/schema"
)

var migrations = []*api.CollectionSchema{
	schema.Sparkle{}.Schema(),
}

func (service *Service) Migrate(drop bool) error {
	ctx := context.TODO()

	resp, err := service.client.Collections().Retrieve(ctx)
	if err != nil {
		return err
	}

	if drop {
		for _, collection := range resp {
			if _, err := service.client.Collection(collection.Name).Delete(ctx); err != nil {
				return err
			}
		}
		resp = nil
	}

	present := make(map[string]*api.CollectionResponse, len(resp))
	for _, collection := range resp {
		present[collection.Name] = collection
	}

	expected := make(map[string]*api.CollectionSchema, len(migrations))
	for _, scheme := range migrations {
		expected[scheme.Name] = scheme
	}

	for name, scheme := range expected {
		if _, ok := present[name]; ok {
			continue
		}
		if _, err := service.client.Collections().Create(ctx, scheme); err != nil {
			return err
		}
	}

	for name := range present {
		if _, ok := expected[name]; ok {
			continue
		}
		if _, err := service.client.Collection(name).Delete(ctx); err != nil {
			return err
		}
	}

	return nil
}
