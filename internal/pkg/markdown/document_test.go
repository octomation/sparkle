package markdown

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacer(t *testing.T) {
	t.Run("replace bytes", func(t *testing.T) {
		var doc Document
		doc.body = []byte("hello, {placeholder}!")

		replace := Replacer(&doc, bytes.ReplaceAll)
		replace([]byte(`{placeholder}`), []byte(`world`))
		assert.Equal(t, []byte("hello, world!"), doc.body)
	})

	t.Run("replace strings", func(t *testing.T) {
		var doc Document
		doc.body = []byte("hello, {placeholder}!")

		replace := Replacer(&doc, strings.ReplaceAll)
		replace("{placeholder}", "world")
		assert.Equal(t, []byte("hello, world!"), doc.body)
	})

	t.Run("replace strings by regexp", func(t *testing.T) {
		var doc Document
		doc.body = []byte("hello, {placeholder}!")
		r := regexp.MustCompile(regexp.QuoteMeta(`{placeholder}`))

		replace := RegexpReplacer(&doc, r.ReplaceAllString)
		replace("world")
		assert.Equal(t, []byte("hello, world!"), doc.body)
	})

	t.Run("replace bytes by regexp", func(t *testing.T) {
		var doc Document
		doc.body = []byte("hello, {placeholder}!")
		r := regexp.MustCompile(regexp.QuoteMeta(`{placeholder}`))

		replace := RegexpReplacer(&doc, r.ReplaceAll)
		replace([]byte(`world`))
		assert.Equal(t, []byte("hello, world!"), doc.body)
	})
}
