package playbook

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	"go.octolab.org/ecosystem/sparkle/internal/pkg/markdown"
	xfs "go.octolab.org/ecosystem/sparkle/internal/pkg/x/fs"
)

const (
	ext = ".md"
	key = "uid"
)

func New(fs afero.Fs) *Service {
	return &Service{fs: fs}
}

type Service struct {
	fs afero.Fs
}

func (service *Service) Notes(root string) ([]Note, error) {
	notes := make([]Note, 0, 100)

	err := afero.Walk(service.fs, root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ext {
			var doc markdown.Document

			file, err := service.fs.Open(path)
			if err != nil {
				return err
			}
			defer safe.Close(file, unsafe.Ignore)

			if err := markdown.LoadFrom(file, &doc); err != nil {
				return err
			}
			notes = append(notes, Note{
				Document:  doc,
				Path:      path,
				CreatedAt: xfs.CreatedAt(info),
				UpdatedAt: info.ModTime(),
			})
		}
		return nil
	})

	return notes, err
}
