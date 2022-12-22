// package fsdd ...
package fsdd

// import
import (
	"time"
)

// const
const (
	_unit     = "bytes"
	_e        = "]"
	_time     = " [time: "
	_dat      = " [data blocks: "
	_inodes   = " [inode(s): "
	_symlinks = " [symlink(s): "
)

// reportInit ...
func reportInit() {
	hashfunc := " [hash:sha512_trunc_256] "
	if c.FastHash {
		hashfunc = " [hash:xxh3] "
	}
	out("FSDD [start] [" + c.Path + "] " + c.Opt + hashfunc)
}

// reportSummary ...
func reportSummary(t0 time.Time) {
	sinodes -= failinodes
	stotal -= failsize
	ssym -= failssym
	if len(c.Opt) > 5 {
		c.Opt = c.Opt[:len(c.Opt)-1]
	}
	prefix := "FSDD [_info] ["
	if !c.HardLink && !c.SymLink && !c.ReplaceSymlinks {
		prefix = "FSDD [_info] [possible "
	}
	ls := len(sym)
	li := len(syi)
	s := " [sym.valid: " + itoa(ls) + "] [sym.invalid: " + itoa(li) + _e
	out("FSDD [_done]" + _time + time.Since(t0).String() + _e)
	out("FSDD [stats] [files:" + itoa(len(sfs)) + _e + _inodes + itoa(len(xfs)+ls+li) + _e + s + _dat + hruIEC(total, _unit) + _e)
	if hcount+hsave > 0 {
		out("FSDD [_info] [deduplication savings]" + _inodes + itoaU64(hcount) + _e + _dat + hruIEC(hsave, _unit) + _e)
	}
	if sinodes+stotal+ssym+uint64(ls) > 0 {
		out(prefix + "new deduplication savings]" + _inodes + itoaU64(sinodes+uint64(len(sym))) + _e + nsyml(c) + _dat + hruIEC(stotal, _unit) + _e)
	}
	if failinodes+failsize+failssym > 0 {
		out("FSDD [_fail] failed to save" + _inodes + itoaU64(failinodes) + itoaU64(failssym) + _symlinks + _dat + hruIEC(failsize, _unit) + _e)
	}
}

func nsyml(c *Config) string {
	if !c.HardLink || !c.SymLink || !c.ReplaceSymlinks {
		return _symlinks + itoa(len(sym)) + _e
	}
	return _symlinks + itoaU64(ssym) + _e
}
