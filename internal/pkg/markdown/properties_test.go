package markdown

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFrontMatter(t *testing.T) {
	var fm FrontMatter
	before := fm.UID()
	assert.NotNil(t, before)
	assert.NotEqual(t, uuid.Nil, before)

	after := uuid.New()
	fm.Set(KeyID, after)
	assert.Equal(t, fm.UID(), fm.Get(KeyID))
	assert.Equal(t, after, fm.UID())
	assert.NotEqual(t, before, after)
}
