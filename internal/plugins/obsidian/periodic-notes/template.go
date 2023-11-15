package periodic

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nleeper/goment"
	xtime "go.octolab.org/time"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

const ext = ".md"

func stub(*goment.Goment) markdown.Transformer {
	return func(*markdown.Document) {}
}

// LinkRelatives replaces a template link with the actual link to the related note.
// E.g., `[[Weekly plans]]` → `[[Week 1, 2006]]`.
func LinkRelatives(cnf Period) func(day *goment.Goment) markdown.Transformer {
	if !cnf.Enabled {
		return stub
	}
	name := strings.TrimSuffix(filepath.Base(cnf.Template), ext)
	if name == "" {
		return stub
	}

	return func(day *goment.Goment) markdown.Transformer {
		old := fmt.Sprintf("[[%s]]", name)
		byNew := fmt.Sprintf("[[%s]]", day.Format(cnf.Format))
		return func(note *markdown.Document) {
			replace := markdown.Replacer(note, strings.ReplaceAll)
			replace(old, byNew)
		}
	}
}

func LinkSiblings(cnf Period) func(day *goment.Goment) markdown.Transformer {
	if !cnf.Enabled {
		return stub
	}

	return func(day *goment.Goment) markdown.Transformer {
		return func(note *markdown.Document) {
			yesterday := (*day).Add(-xtime.Day)
			tomorrow := (*day).Add(+xtime.Day)
			replace := markdown.Replacer(note, strings.ReplaceAll)
			replace("[[prev]]", fmt.Sprintf("[[%s|prev]]", yesterday.Format(cnf.Format)))
			replace("[[next]]", fmt.Sprintf("[[%s|next]]", tomorrow.Format(cnf.Format)))
		}
	}
}
