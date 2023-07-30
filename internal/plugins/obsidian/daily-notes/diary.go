package daily_notes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/nleeper/goment"
	"github.com/spf13/afero"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/errors"
)

const ext = ".md"

func New(cnf Config, opts ...Option) *Diary {
	diary := Diary{
		cnf: cnf,
		ext: ext,
		fs:  afero.NewOsFs(),
	}
	for _, opt := range opts {
		opt(&diary)
	}

	if cnf.Template != "" {
		if filepath.Ext(cnf.Template) == "" {
			cnf.Template += diary.ext
		}
		diary.tpl = Template{
			path: cnf.Template,
			fs:   diary.fs,
		}
	}
	return &diary
}

type Diary struct {
	cnf Config
	ext string
	tpl Template
	fs  afero.Fs
}

func (d *Diary) Create(day time.Time, rewrite bool) (Entry, error) {
	g, err := goment.New(day)
	if err != nil {
		return Entry{}, err // TODO: wrap
	}
	entry := d.entry(g)

	_, err = d.fs.Stat(entry.path)
	exists := !os.IsNotExist(err)
	if exists && !rewrite {
		return entry, nil
	}

	flag := os.O_RDWR | os.O_CREATE
	if rewrite {
		flag |= os.O_TRUNC
	}
	f, err := d.fs.OpenFile(entry.path, flag, 0666)
	if err != nil {
		return entry, err // TODO: wrap
	}
	defer safe.Close(f, unsafe.Ignore)

	structure, err := d.tpl.Structure()
	if err != nil {
		return entry, err
	}
	var note Note
	if err := mapstructure.Decode(structure.FrontMatter, &note.Properties); err != nil {
		return entry, err
	}
	note.Entry = entry
	note.Content = structure.Content
	note.Callbacks = []func(*Note){
		SetUID(uuid.New().String()),
		SetAliases(fmt.Sprintf("Day %d", int(1+day.Sub(d.first().Day()).Hours()/24))),
		SetDate(day.Format(time.DateOnly)),
		LinkWeek(), LinkPrev(), LinkNext(),
	}
	if err := note.Write(f); err != nil {
		return entry, err
	}
	return entry, nil
}

func (d *Diary) entry(g *goment.Goment) Entry {
	entry := Entry{day: g, path: g.Format(d.cnf.Format) + d.ext, format: d.cnf.Format}
	if d.cnf.Folder != "" {
		entry.path = filepath.Join(d.cnf.Folder, entry.path)
	}
	return entry
}

func (d *Diary) first() Entry {
	now, _ := goment.New()
	pattern := filepath.Join(d.cnf.Folder, "*"+d.ext)
	matches, _ := afero.Glob(d.fs, pattern)
	if len(matches) == 0 {
		return d.entry(now)
	}

	// the algorithm provides chronological order of files, not lexicographical
	type name = string
	type path = string
	files := make(map[name]path, len(matches))
	for i := range matches {
		files[strings.TrimSuffix(filepath.Base(matches[i]), d.ext)] = matches[i]
	}

	first, found := new(goment.Goment), ""
	for fname := range files {
		g, err := goment.New(fname, d.cnf.Format)
		if err != nil {
			continue
		}
		if found == "" {
			first = g
			found = fname
			continue
		}
		if g.IsBefore(first) {
			first = g
			found = fname
		}
	}
	if found == "" {
		return d.entry(now)
	}
	return Entry{day: first, path: files[found], format: d.cnf.Format}
}

type Entry struct {
	day    *goment.Goment
	path   string
	format string
}

func (e Entry) Day() time.Time {
	return e.day.ToTime()
}

func (e Entry) Path() string {
	return e.path
}

type Template struct {
	path string
	fs   afero.Fs
}

func (tpl Template) Structure() (pageparser.ContentFrontMatter, error) {
	var stub pageparser.ContentFrontMatter
	f, err := tpl.fs.Open(tpl.path)
	if err != nil {
		return stub, errors.X{
			User:   templateError,
			System: fmt.Errorf("cannot load template %q: %w", tpl.path, err),
		}
	}
	defer safe.Close(f, unsafe.Ignore)

	structure, err := pageparser.ParseFrontMatterAndContent(f)
	if err != nil {
		return stub, errors.X{
			User:   templateError,
			System: fmt.Errorf("cannot parse template %q: %w", tpl.path, err),
		}
	}
	return structure, nil
}
