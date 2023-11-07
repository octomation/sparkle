package periodic

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

// LinkRelatives replaces a template link with the actual link to the related note.
// E.g., `[[Weekly plans]]` → `[[Week 1, 2006]]`.
func LinkRelatives(cnf Period) func(day *goment.Goment) func(*markdown.Document) {
	if !cnf.Enabled {
		return stub
	}
	name := strings.TrimSuffix(filepath.Base(cnf.Template), ext)
	if name == "" {
		return stub
	}
	r, err := regexp.Compile(regexp.QuoteMeta(fmt.Sprintf("[[%s]]", name)))
	if err != nil {
		return stub
	}

	format := cnf.Format
	return func(day *goment.Goment) func(*markdown.Document) {
		return func(note *markdown.Document) {
			note.Content = r.ReplaceAll(
				note.Content,
				[]byte(fmt.Sprintf("[[%s]]", day.Format(format))),
			)
		}
	}
}
