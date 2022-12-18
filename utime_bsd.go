//go:build netbsd || freebsd

// package fsdd ...
package fsdd

// import
import (
	"os"
	"syscall"
)

// rawAtime ...
func rawAtime(fi os.FileInfo) (f, p int64) {
	return fi.Sys().(*syscall.Stat_t).Atimespec.Unix()
}

// rawMtime ...
func rawMtime(fi os.FileInfo) (f, p int64) {
	return fi.Sys().(*syscall.Stat_t).Mtimespec.Unix()
}
