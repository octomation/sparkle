package markdown

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSplitter(t *testing.T) {
	uid := uuid.New()
	fm := []byte(`
uid: {uid}
tags:
  - tag1
  - tag2
`)
	md := []byte(`
---
{fm}
---
# Title

Content.
`)
	fm = bytes.ReplaceAll(fm, []byte("{uid}"), []byte(uid.String()))
	md = bytes.ReplaceAll(md, []byte("{fm}"), fm)

	var splitter Splitter
	props, content, err := splitter.Split(md)
	assert.NoError(t, err)
	assert.Equal(t, string(bytes.TrimSpace(fm)), string(bytes.TrimSpace(props)))
	assert.Equal(t, "# Title\n\nContent.\n", string(content))
}
