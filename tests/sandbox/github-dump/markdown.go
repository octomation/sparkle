package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	InternalID  uuid.UUID `json:"uid" yaml:"uid"`
	ExternalID  []string  `json:"xid" yaml:"xid"`
	Aliases     []string  `json:"aliases" yaml:"aliases"`
	Tags        []string  `json:"tags" yaml:"tags"`
	Topics      []string  `json:"topics" yaml:"topics"`
	Description string    `json:"description" yaml:"description"`
	URL         string    `json:"url" yaml:"url"`
	Homepage    string    `json:"homepage,omitempty" yaml:"homepage,omitempty"`
}

type Markdown struct {
	FrontMatter
	Content string
}

func (md Markdown) DumpTo(file string) error {
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer safe.Close(f, unsafe.Ignore)
	return Encoder{f}.Encode(md)
}

type Encoder struct{ Output io.Writer }

func (e Encoder) Encode(md Markdown) error {
	if _, err := fmt.Fprintln(e.Output, "---"); err != nil {
		return err
	}
	enc := yaml.NewEncoder(e.Output)
	enc.SetIndent(2)
	if err := enc.Encode(md.FrontMatter); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(e.Output, "---"); err != nil {
		return err
	}

	if _, err := fmt.Fprint(e.Output, md.Content); err != nil {
		return err
	}
	return nil
}
