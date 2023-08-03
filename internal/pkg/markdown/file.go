package markdown

import "io"

type File interface {
	io.Reader
	io.Seeker
	io.Writer
	Truncate(int64) error
}

func LoadFrom(file File, doc *Document) error {
	return NewDecoder(file).Decode(doc)
}

func SaveTo(file File, doc Document) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	return NewEncoder(file).Encode(&doc)
}
