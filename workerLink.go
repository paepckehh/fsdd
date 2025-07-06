package fsdd

import (
	"os"
)

//
// WORKER SECTION
//

func feedLinkStats() {
	var result [][]string
	switch c.FastHash {
	case true:
		result = consolidateFastHash()
	default:
		result = consolidateHash()
	}
	for _, v := range result {
		l := len(v)
		var save uint64
		if c.Verbose {
			out("parent [" + v[0] + "]")
		}
		for i := 1; i < l; i++ {
			if ifs[v[0]] != ifs[v[i]] {
				switch {
				case c.HardLink:
					linkChan <- &file{name: v[0], newlinktarget: v[i], size: sfs[v[0]]}
					save += sfs[v[i]]
					sinodes++
				case c.SymLink:
					linkChan <- &file{name: v[0], newlinktarget: v[i], size: sfs[v[0]], newSymLink: true}
					save += sfs[v[i]]
				case c.Verbose:
					out("-> identic [" + v[i] + "]")
					save += sfs[v[i]]
					sinodes++
				default:
					save += sfs[v[i]]
					sinodes++
				}
			}
		}
		stotal += save
	}
	ctl.Done()
}

func workerLink(worker int) {
	for i := 0; i < worker; i++ {
		go func() {
			for f := range linkChan {
				if !linkFast(f) {
					failChan <- f
				}
			}
			bg.Done()
		}()
	}
}

func collectStats() {
	for f := range failChan {
		failsize += f.size
		switch {
		case f.newSymLink:
			failssym++
		case f.syminvalid:
		default:
			failinodes++
		}
	}
	global.Done()
}

func consolidateFastHash() [][]string {
	var r [][]string
	for _, v := range hFs {
		if len(v) > 1 {
			r = append(r, v)
		}
	}
	return r
}

func consolidateHash() [][]string {
	var r [][]string
	for _, v := range hfs {
		if len(v) > 1 {
			r = append(r, v)
		}
	}
	return r
}

//
// BACKEND
//

const _templink = ".fsdd.temp.link.pls.remove.me"

func linkFast(f *file) bool {
	var err error
	if c.Verbose {
		switch {
		case f.newSymLink:
			out("-> new SymLink [" + f.name + "] -> [" + f.newlinktarget + "]")
		default:
			out("-> new HardLink [" + f.name + "] -> [" + f.newlinktarget + "]")
		}
	}
	tlink := f.newlinktarget + _templink
	err = os.Rename(f.newlinktarget, tlink)
	if err != nil {
		errOut("[link] unable to rename [" + f.newlinktarget + "] [" + err.Error() + "]")
		return false
	}
	switch {
	case f.newSymLink:
		err = os.Symlink(f.name, f.newlinktarget)
	default:
		err = os.Link(f.name, f.newlinktarget)
	}
	if err != nil {
		errOut("[link] [unable to link] [" + f.name + "] [" + f.newlinktarget + "] [" + err.Error() + "]")
		errExit("[link] unrecoverable error, please restore [" + f.newlinktarget + "] via [" + f.name + "] manually, EXIT")
	}
	err = os.Remove(tlink)
	if err != nil {
		errOut("[link] unable to rename [" + tlink + "] [" + err.Error() + "]")
		return false
	}
	return true
}
