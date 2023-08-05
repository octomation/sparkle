package schema

import (
	"github.com/google/uuid"
	"github.com/typesense/typesense-go/typesense/api"
	"go.octolab.org/pointer"
)

type Sparkle struct {
	ID         uuid.UUID      `json:"id"`
	Path       string         `json:"path"`
	Properties map[string]any `json:"properties"`
	Content    string         `json:"content"`
	CreatedAt  int64          `json:"created_at"`
	UpdatedAt  int64          `json:"updated_at"`

	Highlights []api.SearchHighlight `json:"-"`
}

func (Sparkle) Schema() *api.CollectionSchema {
	return &api.CollectionSchema{
		Name:                Sparkle{}.Collection(),
		Fields:              Sparkle{}.Fields(),
		DefaultSortingField: pointer.ToString("updated_at"),
		EnableNestedFields:  pointer.ToBool(true),
	}
}

func (Sparkle) Collection() string {
	return "sparkles"
}

func (Sparkle) Fields() []api.Field {
	return []api.Field{
		{
			Name: "id",
			Type: "string",
		},
		{
			Name:  "path",
			Type:  "string",
			Index: pointer.ToBool(true),
		},
		{
			Name: "properties",
			Type: "object",
		},
		{
			Name: "content",
			Type: "string",
		},
		{
			Name:  "created_at",
			Type:  "int64",
			Index: pointer.ToBool(true),
		},
		{
			Name:  "updated_at",
			Type:  "int64",
			Index: pointer.ToBool(true),
		},
	}
}

func (Sparkle) QueryFields() string {
	return "content"
}
