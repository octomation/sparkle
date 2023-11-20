package periodic

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nleeper/goment"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
	xtime "go.octolab.org/ecosystem/sparkle/internal/pkg/x/time"
)

const ext = ".md"

func stub(*goment.Goment) markdown.Transformer {
	return func(*markdown.Document) {}
}

func UpdateAliases() Transformer {
	// based on documentation, https://momentjs.com/docs/#/displaying/format/
	// > To escape characters in format strings, you can wrap
	// > the characters in square brackets.
	// TODO: support similar "MMMM, YYYY" format
	primitive := regexp.MustCompile(`\[\w+]`)
	return func(ref *goment.Goment) markdown.Transformer {
		return func(note *markdown.Document) {
			aliases := note.Property(markdown.KeyAliases)
			if aliases == nil {
				return
			}
			switch aliases := aliases.(type) {
			case []string:
				for i, alias := range aliases {
					if primitive.MatchString(alias) {
						aliases[i] = ref.Format(alias)
					}
				}
			case []any:
				for i, alias := range aliases {
					str, is := alias.(string)
					if !is {
						continue
					}
					if primitive.MatchString(str) {
						aliases[i] = ref.Format(str)
					}
				}
			}
			note.SetProperty(markdown.KeyAliases, aliases)
		}
	}
}

// LinkRelatives replaces a template link with the actual link to the related note.
// E.g., `[[Weekly plans]]` → `[[Week 1, 2006]]`.
func LinkRelatives(cnf Period) Transformer {
	if !cnf.Enabled {
		return stub
	}
	name := strings.TrimSuffix(filepath.Base(cnf.Template), ext)
	if name == "" {
		return stub
	}

	return func(ref *goment.Goment) markdown.Transformer {
		old := fmt.Sprintf("[[%s]]", name)
		byNew := fmt.Sprintf("[[%s]]", ref.Format(cnf.Format))
		return func(note *markdown.Document) {
			replace := markdown.Replacer(note, strings.ReplaceAll)
			replace(old, byNew)
		}
	}
}

func LinkSiblings(
	cnf Period,
	lookup func(goment.Goment) (goment.Goment, goment.Goment),
) Transformer {
	if !cnf.Enabled {
		return stub
	}

	return func(ref *goment.Goment) markdown.Transformer {
		prev, next := lookup(*ref)
		return func(note *markdown.Document) {
			replace := markdown.Replacer(note, strings.ReplaceAll)
			replace("[[prev]]", fmt.Sprintf("[[%s|prev]]", prev.Format(cnf.Format)))
			replace("[[next]]", fmt.Sprintf("[[%s|next]]", next.Format(cnf.Format)))
		}
	}
}

func Lookup(
	prev, next func(time.Time) time.Time,
) func(goment.Goment) (goment.Goment, goment.Goment) {
	return func(ref goment.Goment) (goment.Goment, goment.Goment) {
		t := ref.ToTime()
		pt, nt := prev(t), next(t)
		prev, err := goment.New(pt)
		if err != nil {
			prev = new(goment.Goment)
		}
		next, err := goment.New(nt)
		if err != nil {
			next = new(goment.Goment)
		}
		return *prev, *next
	}
}

func LookupDays(ref goment.Goment) (goment.Goment, goment.Goment) {
	return Lookup(xtime.Yesterday, xtime.Tomorrow)(ref)
}

func LookupWeeks(ref goment.Goment) (goment.Goment, goment.Goment) {
	return Lookup(xtime.PrevWeek, xtime.NextWeek)(ref)
}

func LookupMonths(ref goment.Goment) (goment.Goment, goment.Goment) {
	return Lookup(xtime.PrevMonth, xtime.NextMonth)(ref)
}

func LookupQuarters(ref goment.Goment) (goment.Goment, goment.Goment) {
	return Lookup(xtime.PrevQuarter, xtime.NextQuarter)(ref)
}

func LookupYears(ref goment.Goment) (goment.Goment, goment.Goment) {
	return Lookup(xtime.PrevYear, xtime.NextYear)(ref)
}
