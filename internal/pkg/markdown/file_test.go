package markdown

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {
	const name = "file.md"
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, name, []byte(`---
foo: bar
bar: baz
---
# Baz`), 0666))
	f, err := fs.OpenFile(name, os.O_RDWR, 0666)
	require.NoError(t, err)
	expected := `---
foo: bar
---
# Baz`

	var doc Document
	require.NoError(t, LoadFrom(f, &doc))
	assert.Equal(t, map[string]any{"foo": "bar", "bar": "baz"}, doc.Properties)
	assert.Equal(t, []byte("# Baz"), doc.Content)

	delete(doc.Properties, "bar")
	assert.NoError(t, SaveTo(f, doc))
	assert.NoError(t, f.Close())

	f, err = fs.Open("file.md")
	require.NoError(t, err)
	buf := bytes.NewBuffer(nil)
	_, err = buf.ReadFrom(f)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
	assert.NoError(t, f.Close())
}
