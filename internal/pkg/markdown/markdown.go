package markdown

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gohugoio/hugo/parser/pageparser"
	"gopkg.in/yaml.v3"
)

type Document struct {
	Properties map[string]any
	Content    []byte

	ordered any
}

func (d *Document) SetOrdered(properties any) {
	d.ordered = properties
}

func (d *Document) MarshalMarkdown() ([]byte, error) {
	w := bytes.NewBuffer(nil)

	var props any
	if len(d.Properties) > 0 {
		props = d.Properties
	}
	if d.ordered != nil {
		props = d.ordered
	}
	if props != nil {
		if _, err := fmt.Fprintln(w, "---"); err != nil {
			return nil, err
		}
		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		if err := enc.Encode(props); err != nil {
			return nil, err
		}
		if _, err := fmt.Fprintln(w, "---"); err != nil {
			return nil, err
		}
	}

	if _, err := w.Write(d.Content); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (d *Document) UnmarshalMarkdown(raw []byte) error {
	r := bytes.NewReader(raw)
	structure, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		return err
	}
	d.Properties = structure.FrontMatter
	d.Content = structure.Content
	return nil
}

type Marshaler interface {
	MarshalMarkdown() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalMarkdown([]byte) error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{output: w}
}

type Encoder struct {
	output io.Writer
}

func (e *Encoder) Encode(v any) error {
	raw, err := v.(Marshaler).MarshalMarkdown()
	if err != nil {
		return err
	}
	_, err = e.output.Write(raw)
	return err
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{input: r}
}

type Decoder struct {
	input io.Reader
}

func (d *Decoder) Decode(v any) error {
	raw, err := io.ReadAll(d.input)
	if err != nil {
		return err
	}
	return v.(Unmarshaler).UnmarshalMarkdown(raw)
}
