package markdown

import "io"

type Marshaler interface {
	MarshalMarkdown() ([]byte, error)
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
