package diary

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiary_first(t *testing.T) {
	trick := "2006-02-01"

	t.Run("no files", func(t *testing.T) {
		cnf := Config{
			Folder: "diary",
			Format: "YYYY-DD-MM",
		}
		fs := afero.NewMemMapFs()
		diary := New(cnf, WithSpecifiedFs(fs))

		record := diary.First()
		path := record.Time().Format(trick)
		assert.False(t, record.Time().IsZero())
		assert.Equal(t, fmt.Sprintf("diary/%s.md", path), record.Path)
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
			require.NoError(t, afero.WriteFile(fs, file, []byte{}, 0644))
		}
		diary := New(cnf, WithSpecifiedFs(fs))

		record := diary.First()
		assert.False(t, record.Time().IsZero())
		{ // TODO:debt unexpected behavior of github.com/nleeper/goment
			assert.Equal(t, "diary/2006.md", record.Path)
		}
		// expected behavior, but inverted
		path := record.Time().Format(trick)
		assert.NotEqual(t, fmt.Sprintf("diary/%s.md", path), record.Path)
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
			require.NoError(t, afero.WriteFile(fs, file, []byte{}, 0644))
		}
		diary := New(cnf, WithSpecifiedFs(fs))

		record := diary.First()
		assert.False(t, record.Time().IsZero())
		assert.Equal(t, "diary/2006-31-01.md", record.Path)
	})
}
