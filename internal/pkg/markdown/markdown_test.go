package markdown

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoder_Encode(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	expected := `---
foo: bar
---
# Baz`

	note := Document{
		Properties: map[string]any{"foo": "bar"},
		Content:    []byte("# Baz"),
	}
	assert.NoError(t, NewEncoder(buf).Encode(&note))
	assert.Equal(t, expected, buf.String())
}

func TestDecoder_Decode(t *testing.T) {
	buf := bytes.NewBufferString(`---
foo: bar
---
# Baz`)

	note := Document{}
	assert.NoError(t, NewDecoder(buf).Decode(&note))
	assert.Equal(t, map[string]any{"foo": "bar"}, note.Properties)
	assert.Equal(t, []byte("# Baz"), note.Content)
}
