package diary

import (
	"time"

	"github.com/nleeper/goment"
	"github.com/spf13/afero"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

type Note struct {
	markdown.Document
	Record

	Transformers []markdown.Transformer
}

func (note *Note) SaveTo(file afero.File) error {
	for _, transform := range note.Transformers {
		transform(&note.Document)
	}
	return markdown.SaveTo(file, note.Document)
}

type Record struct {
	Day    goment.Goment
	Path   string
	Format string
}

func (r Record) Time() time.Time {
	return r.Day.ToTime()
}
