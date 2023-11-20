package periodic

import (
	"time"

	"github.com/nleeper/goment"
	"github.com/spf13/afero"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
	xerrors "go.octolab.org/ecosystem/sparkle/internal/pkg/x/errors"
)

func New(cnf Config, opts ...Option) *Planner {
	planner := &Planner{
		cnf: cnf,
		ext: ext,
		fs:  afero.NewOsFs(),
	}
	for _, option := range append(opts, normalize) {
		option(planner)
	}
	return planner
}

type Transformer = func(*goment.Goment) markdown.Transformer

type Planner struct {
	cnf Config
	ext string
	fs  afero.Fs

	changes []Transformer
}

func (p *Planner) Week(ref time.Time, extra ...Transformer) (Note, error) {
	var empty Note
	if err := validate(p.cnf.Weekly); err != nil {
		return empty, err
	}
	g, err := goment.New(ref)
	if err != nil {
		return empty, err
	}
	tpl := newTemplate(p.fs, p.cnf.Weekly, p.ext)
	return tpl.MakeNote(*g, append(p.changes, extra...)...)
}

func (p *Planner) Month(ref time.Time, extra ...Transformer) (Note, error) {
	var empty Note
	if err := validate(p.cnf.Monthly); err != nil {
		return empty, err
	}
	g, err := goment.New(ref)
	if err != nil {
		return empty, err
	}
	tpl := newTemplate(p.fs, p.cnf.Monthly, p.ext)
	return tpl.MakeNote(*g, append(p.changes, extra...)...)
}

func (p *Planner) Quarter(ref time.Time, extra ...Transformer) (Note, error) {
	var empty Note
	if err := validate(p.cnf.Quarterly); err != nil {
		return empty, err
	}
	g, err := goment.New(ref)
	if err != nil {
		return empty, err
	}
	tpl := newTemplate(p.fs, p.cnf.Quarterly, p.ext)
	return tpl.MakeNote(*g, append(p.changes, extra...)...)
}

func (p *Planner) Year(ref time.Time, extra ...Transformer) (Note, error) {
	var empty Note
	if err := validate(p.cnf.Yearly); err != nil {
		return empty, err
	}
	g, err := goment.New(ref)
	if err != nil {
		return empty, err
	}
	tpl := newTemplate(p.fs, p.cnf.Yearly, p.ext)
	return tpl.MakeNote(*g, append(p.changes, extra...)...)
}

func validate(config Period) error {
	if !config.Enabled {
		return xerrors.X{
			User:   errConfig,
			System: nil,
		}
	}
	return nil
}
