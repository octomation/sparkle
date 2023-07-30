package daily_notes

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiary_first(t *testing.T) {
	t.Run("no files", func(t *testing.T) {
		cnf := Config{
			Folder: "diary",
			Format: "YYYY-DD-MM",
		}
		fs := afero.NewMemMapFs()
		diary := New(cnf, WithSpecifiedFs(fs))

		entry := diary.first()
		path := entry.Day().Format("2006-02-01")
		assert.False(t, entry.Day().IsZero())
		assert.Equal(t, fmt.Sprintf("diary/%s.md", path), entry.Path())
	})

	t.Run("no matched files", func(t *testing.T) {
		cnf := Config{
			Folder: "diary",
			Format: "YYYY-DD-MM",
		}
		fs := afero.NewMemMapFs()
		for _, file := range []string{
			"diary/2006-02-01 note.md",
			"diary/2006-02.md",
			"diary/2006.md",
		} {
			require.NoError(t, afero.WriteFile(fs, file, []byte{}, 0666))
		}
		diary := New(cnf, WithSpecifiedFs(fs))

		entry := diary.first()
		assert.False(t, entry.Day().IsZero())
		{ // TODO:debt unexpected behavior of github.com/nleeper/goment
			assert.Equal(t, "diary/2006.md", entry.Path())
		}
		// expected behavior, but inverted
		path := entry.Day().Format("2006-02-01")
		assert.NotEqual(t, fmt.Sprintf("diary/%s.md", path), entry.Path())
	})

	t.Run("chronological order", func(t *testing.T) {
		cnf := Config{
			Folder: "diary",
			Format: "YYYY-DD-MM",
		}
		fs := afero.NewMemMapFs()
		for _, file := range []string{
			"diary/2006-01-02.md",
			"diary/2006-02-02.md",
			"diary/2006-03-02.md",
			"diary/2006-31-01.md",
		} {
			require.NoError(t, afero.WriteFile(fs, file, []byte{}, 0666))
		}
		diary := New(cnf, WithSpecifiedFs(fs))

		entry := diary.first()
		assert.False(t, entry.Day().IsZero())
		assert.Equal(t, "diary/2006-31-01.md", entry.Path())
	})
}
