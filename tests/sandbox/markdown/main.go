package main

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/russross/blackfriday/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func main() {
	raw := []byte(`
---
foo: bar
bar: baz
---
# Hello!
`)
	r := text.NewReader(raw)
	p := goldmark.DefaultParser()
	n := p.Parse(r)
	n.Dump(raw, 0)

	md := blackfriday.New()
	x := md.Parse(raw)
	fmt.Println(x.String())

	raw = []byte(`
foo: bar
bar: baz
`)
	var doc any
	dec := yaml.NewDecoder(bytes.NewReader(raw), yaml.UseOrderedMap())
	fmt.Println(dec.Decode(&doc))
	fmt.Println(doc)
}
