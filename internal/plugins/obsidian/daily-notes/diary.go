package diary

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nleeper/goment"
	"github.com/spf13/afero"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/errors"
	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
)

const ext = ".md"

func New(cnf Config, opts ...Option) *Diary {
	diary := &Diary{
		cnf: cnf,
		ext: ext,
		fs:  afero.NewOsFs(),
	}
	for _, opt := range append(opts, normalize) {
		opt(diary)
	}
	return diary
}

type Diary struct {
	cnf Config
	ext string
	fs  afero.Fs
}

func (d *Diary) Create(
	day time.Time,
	transformers ...func(*goment.Goment) func(*markdown.Document),
) (Record, error) {
	// TODO:feat support --rewrite, --merge, --ignore strategies and `unknown` err by default
	rewrite := true

	g, _ := goment.New(day)
	record := d.record(g)
	_, err := d.fs.Stat(record.Path)
	if err == nil && !rewrite {
		return record, nil
	}

	flag := os.O_RDWR | os.O_CREATE
	if rewrite {
		flag |= os.O_TRUNC
	}
	file, err := d.fs.OpenFile(record.Path, flag, 0644)
	if err != nil {
		return record, errors.X{
			User:   errFolder,
			System: fmt.Errorf("cannot open file %q: %w", record.Path, err),
		}
	}
	defer safe.Close(file, unsafe.Ignore)

	template, err := d.template()
	if err != nil {
		return record, err
	}

	note := Note{Document: template, Record: record}
	note.Transformers = []func(*Note){
		SetUID(uuid.New().String()),
		SetAliases(
			fmt.Sprintf(
				"Day %d",
				int(1+day.Sub(d.First().Time()).Hours()/24),
			),
		),
		SetDate(day.Format(time.DateOnly)),
		LinkPrev(), LinkNext(),
	}

	// TODO:deps:refactoring improve configuration
	for _, fn := range transformers {
		note.TransformersNew = append(note.TransformersNew, fn(g))
	}

	if err := note.SaveTo(file); err != nil {
		return record, errors.X{
			User:   errFolder,
			System: fmt.Errorf("cannot save file %q: %w", record.Path, err),
		}
	}
	return record, nil
}

func (d *Diary) First() Record {
	return d.find(func(sup, iter *goment.Goment) bool { return iter.IsBefore(sup) })
}

func (d *Diary) Last() Record {
	return d.find(func(sup, iter *goment.Goment) bool { return iter.IsAfter(sup) })
}

func (d *Diary) find(cmp func(sup, iter *goment.Goment) bool) Record {
	now, _ := goment.New()
	pattern := filepath.Join(d.cnf.Folder, "*"+d.ext)
	matches, _ := afero.Glob(d.fs, pattern)
	if len(matches) == 0 {
		return d.record(now)
	}

	// the algorithm provides chronological order of files, not lexicographical
	type name = string
	type path = string
	files := make(map[name]path, len(matches))
	for i := range matches {
		files[strings.TrimSuffix(filepath.Base(matches[i]), d.ext)] = matches[i]
	}

	day, found := new(goment.Goment), ""
	for fname := range files {
		g, err := goment.New(fname, d.cnf.Format)
		if err != nil {
			continue
		}
		if found == "" {
			day = g
			found = fname
			continue
		}
		if cmp(day, g) {
			day = g
			found = fname
		}
	}
	if found == "" {
		return d.record(now)
	}
	return Record{Day: *day, Path: files[found], Format: d.cnf.Format}
}

func (d *Diary) record(g *goment.Goment) Record {
	record := Record{Day: *g, Path: g.Format(d.cnf.Format) + d.ext, Format: d.cnf.Format}
	if d.cnf.Folder != "" {
		record.Path = filepath.Join(d.cnf.Folder, record.Path)
	}
	return record
}

func (d *Diary) template() (markdown.Document, error) {
	empty := markdown.Document{
		Properties: make(map[string]any),
		Content:    nil,
	}
	if d.cnf.Template == "" {
		return empty, nil
	}

	file, err := d.fs.Open(d.cnf.Template)
	if err != nil {
		return empty, errors.X{
			User:   errTemplate,
			System: fmt.Errorf("cannot load template %q: %w", d.cnf.Template, err),
		}
	}
	defer safe.Close(file, unsafe.Ignore)

	var doc markdown.Document
	if err := markdown.LoadFrom(file, &doc); err != nil {
		return empty, errors.X{
			User:   errTemplate,
			System: fmt.Errorf("cannot parse template %q: %w", d.cnf.Template, err),
		}
	}
	return doc, nil
}
