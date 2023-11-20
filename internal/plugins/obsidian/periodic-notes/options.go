package periodic

import (
	"path/filepath"

	"github.com/spf13/afero"
)

type Option func(*Planner)

func WithSpecifiedExt(ext string) Option {
	return func(p *Planner) {
		p.ext = ext
	}
}

func WithSpecifiedFs(fs afero.Fs) Option {
	return func(p *Planner) {
		p.fs = fs
	}
}

func WithTransformers(transformers ...Transformer) Option {
	return func(p *Planner) {
		p.changes = transformers
	}
}

func normalize(p *Planner) {
	for _, period := range []*Period{
		&p.cnf.Weekly,
		&p.cnf.Monthly,
		&p.cnf.Quarterly,
		&p.cnf.Yearly,
	} {
		if path := period.Template; path != "" {
			if filepath.Ext(path) == "" {
				period.Template += p.ext
			}
		}
	}
}
