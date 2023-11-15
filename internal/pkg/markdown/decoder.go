package markdown

import "io"

type Unmarshaler interface {
	UnmarshalMarkdown([]byte) error
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
