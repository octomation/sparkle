package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	xtesting "go.octolab.org/ecosystem/sparkle/internal/pkg/x/testing"
)

func TestServerSerialization(t *testing.T) {
	const license = "65541023-fb3e-4107-ac8e-158fc2e64a18"

	f, err := os.Open("testdata/server.toml")
	require.NoError(t, err)
	defer xtesting.Close(t, f)

	var cnf Server
	dec := toml.NewDecoder(f)
	_, err = dec.Decode(&cnf)
	assert.NoError(t, err)
	assert.Equal(t, uuid.MustParse(license), cnf.Service.License)
	assert.FileExists(t, filepath.Join(cnf.Sparkle.Path, cnf.Sparkle.File))
}
