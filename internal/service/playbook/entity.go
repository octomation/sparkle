package playbook

import (
	"time"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

type Note struct {
	markdown.Document
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
