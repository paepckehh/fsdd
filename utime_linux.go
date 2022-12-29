//go:build openbsd || aix || linux

package fsdd

import (
	"os"
	"syscall"
)

func rawAtime(fi os.FileInfo) (f, p int64) {
	return fi.Sys().(*syscall.Stat_t).Atim.Unix()
}

func rawMtime(fi os.FileInfo) (f, p int64) {
	return fi.Sys().(*syscall.Stat_t).Mtim.Unix()
}

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
