package markdown

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	uid := uuid.New()
	md := []byte(`
---
uid: {uid}
tags:
  - tag1
  - tag2
---
# Title

Content.
`)
	md = bytes.ReplaceAll(md, []byte("{uid}"), []byte(uid.String()))

	var doc Document
	dec := NewDecoder(bytes.NewBuffer(md))
	assert.NoError(t, dec.Decode(&doc))
	assert.Equal(t, uid, doc.head.UID())
	assert.Len(t, doc.head.Get(KeyTags), 2)
}
