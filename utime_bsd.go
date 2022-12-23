//go:build netbsd || freebsd

package fsdd

import (
	"os"
	"syscall"
)

func rawAtime(fi os.FileInfo) (f, p int64) {
	return fi.Sys().(*syscall.Stat_t).Atimespec.Unix()
}

func rawMtime(fi os.FileInfo) (f, p int64) {
	return fi.Sys().(*syscall.Stat_t).Mtimespec.Unix()
}
