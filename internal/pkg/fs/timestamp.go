//go:build !linux && !unix && !darwin && !freebsd

package fs

import (
	"io/fs"
	"time"
)

func CreatedAt(info fs.FileInfo) time.Time {
	return info.ModTime()
}
