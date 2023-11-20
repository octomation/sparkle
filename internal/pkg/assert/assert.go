package xassert

import (
	"strings"

	"github.com/pkg/errors"
)

var disabled bool

func True(check func() bool, messages ...string) {
	if disabled {
		return
	}
	if check() {
		return
	}

	message := "assertion failed"
	if len(messages) > 0 {
		message = messages[0]
		if len(messages) > 1 {
			message = strings.Join(messages, " ")
		}
	}
	// I need stack trace here
	panic(errors.New(message))
}
