package daily_notes

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Run("no configuration file", func(t *testing.T) {
		expected := defaults
		fs := afero.NewMemMapFs()

		cnf, err := LoadConfig(fs)
		assert.NoError(t, err)
		assert.Equal(t, expected, cnf)
	})

	t.Run("success loading", func(t *testing.T) {
		expected := Config{
			Autorun:  true,
			Folder:   "path/to/daily/notes",
			Format:   "DD.MM.YYYY",
			Template: "path/to/template.md",
		}
		fs := afero.NewMemMapFs()
		r := bytes.NewBuffer(nil)
		require.NoError(t, json.NewEncoder(r).Encode(expected))
		require.NoError(t, afero.WriteReader(fs, config, r))

		cnf, err := LoadConfig(fs)
		assert.NoError(t, err)
		assert.Equal(t, expected, cnf)
	})

	t.Run("load invalid file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fs, config, []byte("broken"), 0o000))

		_, err := LoadConfig(fs)
		assert.Error(t, err)
	})
}
