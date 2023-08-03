package diary

import (
	"path/filepath"

	"github.com/spf13/afero"
)

type Option func(*Diary)

func normalize(d *Diary) {
	if path := d.cnf.Template; path != "" {
		if filepath.Ext(path) == "" {
			d.cnf.Template += d.ext
		}
	}
}

func WithSpecifiedExt(ext string) Option {
	return func(d *Diary) {
		d.ext = ext
	}
}

func WithSpecifiedFs(fs afero.Fs) Option {
	return func(d *Diary) {
		d.fs = fs
	}
}
