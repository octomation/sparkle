package config

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	xtesting "go.octolab.org/ecosystem/sparkle/internal/pkg/x/testing"
)

func TestSparkleSerialization(t *testing.T) {
	f, err := os.Open("testdata/sparkle.toml")
	require.NoError(t, err)
	defer xtesting.Close(t, f)

	var cnf Sparkle
	dec := toml.NewDecoder(f)
	_, err = dec.Decode(&cnf)
	assert.NoError(t, err)
	assert.True(t, cnf.Obsidian.Plugins.Calendar.Enabled)
	assert.True(t, cnf.Obsidian.Plugins.DailyNotes.Enabled)
	assert.True(t, cnf.Obsidian.Plugins.PeriodicNotes.Enabled)
}
