package markdown

import (
	"bytes"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

const (
	KeyAliases   = "aliases"
	KeyCreatedAt = "created at"
	KeyDate      = "date"
	KeyID        = "uid"
	KeyRemindMe  = "remind me"
	KeyTags      = "tags"
)

type FrontMatter struct {
	props yaml.MapSlice
	dirty bool
}

func (fm *FrontMatter) UID() uuid.UUID {
	return fm.handleID().Get(KeyID).(uuid.UUID)
}

func (fm *FrontMatter) Get(key string) any {
	for _, item := range fm.props {
		if item.Key == key {
			return item.Value
		}
	}
	return nil
}

func (fm *FrontMatter) Set(key string, value any) {
	for i, item := range fm.props {
		if item.Key == key {
			fm.props[i].Value = value
			fm.dirty = true
			return
		}
	}
	fm.props = append(fm.props, yaml.MapItem{Key: key, Value: value})
	fm.dirty = true
}

func (fm *FrontMatter) Delete(key string) {
	for i, item := range fm.props {
		if item.Key == key {
			fm.props = append(fm.props[:i], fm.props[i+1:]...)
			fm.dirty = true
			return
		}
	}
}

func (fm *FrontMatter) IsDirty() bool {
	return fm.dirty
}

func (fm *FrontMatter) UnmarshalYAML(val []byte) error {
	dec := yaml.NewDecoder(bytes.NewBuffer(val))
	err := dec.Decode(&fm.props)
	if err == nil {
		fm.handleID()
	}
	return err
}

func (fm *FrontMatter) MarshalYAML() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buf, yaml.Indent(2), yaml.IndentSequence(true))
	err := enc.Encode(fm.handleID().props)
	return buf.Bytes(), err
}

func (fm *FrontMatter) handleID() *FrontMatter {
	switch id := fm.Get(KeyID).(type) {
	case string:
		uid, err := uuid.Parse(id)
		if err != nil || uid == uuid.Nil {
			uid = uuid.New()
		}
		fm.Set(KeyID, uid)
	case uuid.UUID:
		if id == uuid.Nil {
			fm.Set(KeyID, uuid.New())
		}
	default:
		fm.Set(KeyID, uuid.New())
	}
	return fm
}
