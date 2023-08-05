//go:build darwin || freebsd

package fs

import (
	"io/fs"
	"syscall"
	"time"
)

func CreatedAt(info fs.FileInfo) time.Time {
	sys, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return info.ModTime()
	}
	return time.Unix(sys.Birthtimespec.Unix())
}
