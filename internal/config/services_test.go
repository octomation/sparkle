package config

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToml(t *testing.T) {
	f, err := os.Open("testdata/sparkle.toml")
	require.NoError(t, err)

	var cfg Server
	_, err = toml.NewDecoder(f).Decode(&cfg)
	assert.NoError(t, err)
	assert.Equal(t, "locale", cfg.Service.Plugins.Obsidian.Calendar.WeekStart)
}
