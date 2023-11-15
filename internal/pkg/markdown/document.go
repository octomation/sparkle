package markdown

import (
	"bytes"

	"github.com/google/uuid"
	"go.octolab.org/unsafe"
)

type Document struct {
	head FrontMatter
	body []byte
}

func (doc *Document) ID() uuid.UUID {
	return doc.head.UID()
}

func (doc *Document) Content() string {
	return string(doc.body)
}

func (doc *Document) Properties() map[string]any {
	props := make(map[string]any, len(doc.head.props))
	for _, item := range doc.head.props {
		props[item.Key.(string)] = item.Value
	}
	return props
}

func (doc *Document) Property(key string) any {
	return doc.head.Get(key)
}

func (doc *Document) SetProperty(key string, value any) *Document {
	doc.head.Set(key, value)
	return doc
}

func (doc *Document) DeleteProperty(key string) *Document {
	doc.head.Delete(key)
	return doc
}

func (doc *Document) MarshalMarkdown() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	unsafe.DoSilent(buf.WriteString("---\n"))
	raw, err := doc.head.MarshalYAML()
	if err != nil {
		return nil, err
	}
	unsafe.DoSilent(buf.Write(raw))
	unsafe.DoSilent(buf.WriteString("---\n"))
	unsafe.DoSilent(buf.Write(doc.body))
	return buf.Bytes(), nil
}

func (doc *Document) UnmarshalMarkdown(raw []byte) error {
	splitter := new(Splitter)
	head, body, err := splitter.Split(raw)
	if err != nil {
		return err
	}
	if err := doc.head.UnmarshalYAML(head); err != nil {
		return err
	}
	doc.body = body
	return nil
}

type Content interface {
	[]byte | string
}

func Replacer[T Content](doc *Document, fn func(src, old, new T) T) func(old, new T) {
	return func(old, new T) {
		doc.body = []byte(fn(T(doc.body), old, new))
	}
}

func RegexpReplacer[T Content](doc *Document, fn func(src, repl T) T) func(repl T) {
	return func(repl T) {
		doc.body = []byte(fn(T(doc.body), repl))
	}
}

type Transformer func(*Document)

func RegenerateID(doc *Document) {
	doc.head.Set(KeyID, uuid.New())
}

func AddAliases(aliases ...string) Transformer {
	return func(doc *Document) {
		switch stored := doc.head.Get(KeyAliases).(type) {
		case []string:
			doc.head.Set(KeyAliases, append(stored, aliases...))
		case []any:
			for _, alias := range aliases {
				stored = append(stored, any(alias))
			}
			doc.head.Set(KeyAliases, stored)
		}
	}
}

func SetAliases(aliases ...string) Transformer {
	return func(doc *Document) {
		doc.head.Set(KeyAliases, aliases)
	}
}

func AddTags(tags ...string) Transformer {
	return func(doc *Document) {
		switch stored := doc.head.Get(KeyTags).(type) {
		case []string:
			doc.head.Set(KeyTags, append(stored, tags...))
		case []any:
			for _, alias := range tags {
				stored = append(stored, any(alias))
			}
			doc.head.Set(KeyTags, stored)
		}
	}
}

func SetTags(tags ...string) Transformer {
	return func(doc *Document) {
		doc.head.Set(KeyTags, tags)
	}
}

func SetDate(date string) Transformer {
	return func(doc *Document) {
		doc.head.Set(KeyDate, date)
	}
}
