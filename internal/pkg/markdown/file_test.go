package markdown

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.octolab.org/unsafe"
)

func TestFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	f, err := fs.Create("test.md")
	require.NoError(t, err)

	var doc Document
	reminder := strings.TrimSuffix(time.RFC3339, "Z07:00")
	doc.head.Set(KeyCreatedAt, time.RFC3339)
	doc.head.Set(KeyAliases, []string{"alias1"})
	doc.head.Set(KeyTags, []string{"tag1", "tag2"})
	doc.head.Set(KeyRemindMe, reminder)

	assert.NoError(t, SaveTo(f, doc))
	unsafe.DoSilent(f.Seek(0, io.SeekStart))
	assert.NoError(t, LoadFrom(f, &doc))
	assert.Len(t, doc.head.Get(KeyAliases), 1)
	assert.Len(t, doc.head.Get(KeyTags), 2)
	assert.Equal(t, time.RFC3339, doc.head.Get(KeyCreatedAt))
	assert.Equal(t, reminder, doc.head.Get(KeyRemindMe))
}
