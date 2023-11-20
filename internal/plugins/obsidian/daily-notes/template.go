package diary

import (
	"fmt"
	"strings"

	"github.com/nleeper/goment"
	xtime "go.octolab.org/time"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

const ext = ".md"

func stub(*goment.Goment) markdown.Transformer {
	return func(*markdown.Document) {}
}

func LinkSiblings(cnf Config) func(day *goment.Goment) markdown.Transformer {
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
