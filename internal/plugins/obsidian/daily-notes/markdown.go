package daily_notes

import (
	"fmt"
	"io"
	"regexp"

	"gopkg.in/yaml.v3"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

type Properties struct {
	markdown.Core        `mapstructure:",squash" yaml:",inline"`
	markdown.DiaryEntry  `mapstructure:",squash" yaml:",inline"`
	markdown.MoodTracker `mapstructure:",squash" yaml:",inline"`
	markdown.TimeTracker `mapstructure:",squash" yaml:",inline"`
}

type Note struct {
	Entry
	Properties
	Content   []byte
	Callbacks []func(*Note)
}

func (note *Note) Write(w io.Writer) error {
	for _, fn := range note.Callbacks {
		fn(note)
	}

	// front matter
	if _, err := fmt.Fprintln(w, "---"); err != nil {
		return err
	}
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(note.Properties); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "---"); err != nil {
		return err
	}

	// content
	if _, err := w.Write(note.Content); err != nil {
		return err
	}
	return nil
}

func SetUID(uid string) func(*Note) {
	return func(note *Note) {
		note.UID = uid
	}
}

func SetAliases(aliases ...string) func(*Note) {
	return func(note *Note) {
		note.Aliases = aliases
	}
}

func AddAliases(aliases ...string) func(*Note) {
	return func(note *Note) {
		note.Aliases = append(note.Aliases, aliases...)
	}
}

func AddTags(tags ...string) func(*Note) {
	return func(note *Note) {
		note.Tags = append(note.Tags, tags...)
	}
}

func SetTags(tags ...string) func(*Note) {
	return func(note *Note) {
		note.Tags = tags
	}
}

func SetDate(date string) func(*Note) {
	return func(note *Note) {
		note.Date = date
	}
}

func LinkWeek() func(*Note) {
	grep := regexp.MustCompile(`\[\[Weekly plans]]`)
	return func(note *Note) {
		year, week := note.day.ToTime().ISOWeek()
		replacement := fmt.Sprintf("[[Week %d, %d]]", week, year)
		note.Content = grep.ReplaceAll(note.Content, []byte(replacement))
	}
}

func LinkPrev() func(*Note) {
	grep := regexp.MustCompile(`\[\[prev]]`)
	return func(note *Note) {
		replacement := fmt.Sprintf("[[%s|prev]]", note.day.Add(-1, "d").Format(note.format))
		note.Content = grep.ReplaceAll(note.Content, []byte(replacement))
	}
}

func LinkNext() func(*Note) {
	grep := regexp.MustCompile(`\[\[next]]`)
	return func(note *Note) {
		replacement := fmt.Sprintf("[[%s|next]]", note.day.Add(+1, "d").Format(note.format))
		note.Content = grep.ReplaceAll(note.Content, []byte(replacement))
	}
}
