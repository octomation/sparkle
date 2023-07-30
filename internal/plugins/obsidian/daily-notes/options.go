package daily_notes

import "github.com/spf13/afero"

type Option func(*Diary)

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
