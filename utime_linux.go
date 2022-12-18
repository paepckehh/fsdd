//go:build openbsd || aix || linux

// package fsdd ...
package fsdd

// import
import (
	"os"
	"syscall"
)

// rawAtime ...
func rawAtime(fi os.FileInfo) (f, p int64) {
	return fi.Sys().(*syscall.Stat_t).Atim.Unix()
}

// rawMtime ...
func rawMtime(fi os.FileInfo) (f, p int64) {
	return fi.Sys().(*syscall.Stat_t).Mtim.Unix()
}

// isInodeTimeZero ...
func isInodeTimeZero(fi os.FileInfo) bool {
	stat := fi.Sys().(*syscall.Stat_t)
	as, ans := stat.Atim.Unix()
	ms, mns := stat.Mtim.Unix()
	cs, cns := stat.Ctim.Unix()
	if as+ans+ms+mns+cs+cns > 0 {
		return false
	}
	return true
}
