package calendar

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nleeper/goment"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

const ext = ".md"

func stub(*goment.Goment) func(*markdown.Document) {
	return func(*markdown.Document) {}
}

// LinkWeek replaces a template link with the actual link to the weekly note.
// E.g., `[[Weekly plans]]` → `[[Week 1, 2006]]`.
func LinkWeek(cnf Config) func(day *goment.Goment) func(*markdown.Document) {
	name := strings.TrimSuffix(filepath.Base(cnf.WeeklyNoteTemplate), ext)
	if name == "" {
		return stub
	}
	r, err := regexp.Compile(regexp.QuoteMeta(fmt.Sprintf("[[%s]]", name)))
	if err != nil {
		return stub
	}

	format := cnf.WeeklyNoteFormat
	return func(day *goment.Goment) func(*markdown.Document) {
		return func(note *markdown.Document) {
			note.Content = r.ReplaceAll(
				note.Content,
				[]byte(fmt.Sprintf("[[%s]]", day.Format(format))),
			)
		}
	}
}
