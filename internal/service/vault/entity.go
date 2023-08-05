package vault

import (
	"time"

	"github.com/google/uuid"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

type Note struct {
	markdown.Document
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (note Note) ID() uuid.UUID {
	uid := note.Properties[key]
	if uid == nil {
		return uuid.Nil
	}
	return uuid.MustParse(uid.(string))
}

func (note Note) Content() string {
	return string(note.Document.Content)
}
