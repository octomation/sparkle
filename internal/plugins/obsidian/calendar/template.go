package calendar

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nleeper/goment"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

const ext = ".md"

func stub(*goment.Goment) markdown.Transformer {
	return func(*markdown.Document) {}
}

// LinkWeek replaces a template link with the actual link to the weekly note.
// E.g., `[[Weekly plans]]` → `[[Week 1, 2006]]`.
func LinkWeek(cnf Config) func(day *goment.Goment) markdown.Transformer {
	if !cnf.Enabled {
		return stub
	}
	name := strings.TrimSuffix(filepath.Base(cnf.WeeklyNoteTemplate), ext)
	if name == "" {
		return stub
	}

	return func(day *goment.Goment) markdown.Transformer {
		old := fmt.Sprintf("[[%s]]", name)
		byNew := fmt.Sprintf("[[%s]]", day.Format(cnf.WeeklyNoteFormat))
		return func(note *markdown.Document) {
			replace := markdown.Replacer(note, strings.ReplaceAll)
			replace(old, byNew)
		}
	}
}
