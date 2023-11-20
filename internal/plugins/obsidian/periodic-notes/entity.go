package periodic

import (
	"fmt"
	"path/filepath"

	"github.com/nleeper/goment"
	"github.com/spf13/afero"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
	xerrors "go.octolab.org/ecosystem/sparkle/internal/pkg/x/errors"
)

type Note struct {
	markdown.Document
	Ref goment.Goment
	Src string
}

func newTemplate(fs afero.Fs, cnf Period, ext string) template {
	return template{
		fs:  fs,
		ext: ext,
		src: cnf.Template,
		dst: cnf.Folder,
		fmt: cnf.Format,
	}
}

type template struct {
	// storage
	fs afero.Fs

	// config
	ext string
	src string
	dst string
	fmt string
}

func (tpl template) load() (markdown.Document, error) {
	var doc markdown.Document
	if tpl.src == "" {
		return doc, nil
	}

	f, err := tpl.fs.Open(tpl.src)
	if err != nil {
		return doc, xerrors.X{
			User:   errTemplate,
			System: fmt.Errorf("cannot load template %q: %w", tpl.src, err),
		}
	}
	defer safe.Close(f, unsafe.Ignore)

	if err := markdown.LoadFrom(f, &doc); err != nil {
		return doc, xerrors.X{
			User:   errTemplate,
			System: fmt.Errorf("cannot parse template %q: %w", tpl.src, err),
		}
	}
	return doc, nil
}

func (tpl template) MakeNote(ref goment.Goment, transformers ...Transformer) (Note, error) {
	note := Note{Ref: ref, Src: filepath.Join(tpl.dst, ref.Format(tpl.fmt)+tpl.ext)}
	doc, err := tpl.load()
	if err != nil {
		return note, err
	}
	note.Document = doc

	// TODO:implement rewrite, replace, ignore strategies
	f, err := tpl.fs.Create(note.Src)
	if err != nil {
		return note, xerrors.X{
			User:   errFolder,
			System: fmt.Errorf("cannot create note %q: %w", note.Src, err),
		}
	}
	defer safe.Close(f, unsafe.Ignore)

	for _, transformer := range transformers {
		transform := transformer(&ref)
		transform(&note.Document)
	}
	return note, markdown.SaveTo(f, note.Document)
}
