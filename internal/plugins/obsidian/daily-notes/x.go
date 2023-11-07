package diary

import (
	"fmt"
	"regexp"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/afero"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

type Note struct {
	markdown.Document
	Record

	Transformers    []func(*Note)
	TransformersNew []func(*markdown.Document)
	ordered         struct {
		markdown.Core        `mapstructure:",squash" yaml:",inline"`
		markdown.DiaryEntry  `mapstructure:",squash" yaml:",inline"`
		markdown.MoodTracker `mapstructure:",squash" yaml:",inline"`
		markdown.TimeTracker `mapstructure:",squash" yaml:",inline"`
		Other                map[string]any `mapstructure:",remain" yaml:",inline"`
	}
}

func (note *Note) SaveTo(file afero.File) error {
	if err := mapstructure.Decode(note.Properties, &note.ordered); err != nil {
		return err
	}
	for _, transform := range note.Transformers {
		transform(note)
	}

	// TODO:deps:refactoring improve configuration
	for _, transform := range note.TransformersNew {
		transform(&note.Document)
	}

	note.SetOrdered(note.ordered)

	return markdown.SaveTo(file, note.Document)
}

func SetUID(uid string) func(*Note) {
	return func(note *Note) {
		note.ordered.UID = uid
	}
}

func SetAliases(aliases ...string) func(*Note) {
	return func(note *Note) {
		note.ordered.Aliases = aliases
	}
}

func AddAliases(aliases ...string) func(*Note) {
	return func(note *Note) {
		note.ordered.Aliases = append(note.ordered.Aliases, aliases...)
	}
}

func AddTags(tags ...string) func(*Note) {
	return func(note *Note) {
		note.ordered.Tags = append(note.ordered.Tags, tags...)
	}
}

func SetTags(tags ...string) func(*Note) {
	return func(note *Note) {
		note.ordered.Tags = tags
	}
}

func SetDate(date string) func(*Note) {
	return func(note *Note) {
		note.ordered.Date = date
	}
}

func LinkPrev() func(*Note) {
	grep := regexp.MustCompile(`\[\[prev]]`)
	return func(note *Note) {
		replacement := fmt.Sprintf("[[%s|prev]]", note.Yesterday().Format(note.Format))
		note.Content = grep.ReplaceAll(note.Content, []byte(replacement))
	}
}

func LinkNext() func(*Note) {
	grep := regexp.MustCompile(`\[\[next]]`)
	return func(note *Note) {
		replacement := fmt.Sprintf("[[%s|next]]", note.Tomorrow().Format(note.Format))
		note.Content = grep.ReplaceAll(note.Content, []byte(replacement))
	}
}
