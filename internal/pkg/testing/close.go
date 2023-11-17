package xtesting

import (
	"io"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Tester interface {
	Errorf(string, ...any)
	FailNow()
}

func Close(tester Tester, closer io.Closer) {
	assert.NoError(tester, closer.Close())
}

func StrictClose(tester Tester, closer io.Closer) {
	require.NoError(tester, closer.Close())
}
