package markdown

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEncoder(t *testing.T) {
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
	doc.head.Set(KeyID, uid)
	doc.head.Set(KeyTags, []string{"tag1", "tag2"})
	doc.body = []byte("# Title\n\nContent.")

	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	assert.NoError(t, enc.Encode(&doc))
	assert.Equal(t, string(bytes.TrimSpace(md)), buf.String())
}
