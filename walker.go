// package fsdd ...
package fsdd

// import
import (
	"os"
	"syscall"
)

//
// WORKER SECTION
//

// const
const (
	_modeDir     uint32 = 1 << (32 - 1 - 0)
	_modeSymlink uint32 = 1 << (32 - 1 - 4)
)

// workerTreeWalker ...
func workerTreeWalker(worker int) {
	fi, err := os.Stat(c.Path)
	if err != nil {
		errExit("[stat deviceID root dir] [" + c.Path + "] [" + err.Error() + "]")
	}
	d, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		errExit("[stat deviceID root dir] [" + c.Path + "]")
	}
	rootNodeDeviceID := uint64(d.Dev)
	for i := 0; i < worker; i++ {
		go func() {
			for path := range dirChan {
				list, err := os.ReadDir(path)
				if err != nil {
					errOut("[read dir] [" + path + "] [" + err.Error() + "]")
					return
				}
				for _, item := range list {
					target := path + "/" + item.Name()
					fi, _ := item.Info()
					size := fi.Size()
					switch {
					case item.Type().IsRegular():
						if c.CleanMetadata {
							as, ans := rawAtime(fi)
							ms, mns := rawMtime(fi)
							if as+ans+ms+mns != 0 {
								metaChan <- target
							}
						}
						fileChan <- &file{
							name:  target,
							size:  uint64(size),
							inode: fi.Sys().(*syscall.Stat_t).Ino,
						}
					case uint32(item.Type())&_modeDir != 0:
						if c.CleanMetadata {
							as, ans := rawAtime(fi)
							ms, mns := rawMtime(fi)
							if as+ans+ms+mns != 0 {
								metaChan <- target
							}
						}
						st, _ := fi.Sys().(*syscall.Stat_t)
						if uint64(st.Dev) != rootNodeDeviceID {
							continue // skip dir targets outside of our fs boundary
						}
						bg.Add(1)
						dirChan <- target
					case uint32(item.Type())&_modeSymlink != 0:
						symtarget, err := os.Readlink(target)
						if err != nil {
							fileChan <- &file{symlink: true, syminvalid: true, name: target}
							out("unable to read symlink [" + target + "] -> [" + err.Error() + "]")
							continue
						}
						switch symtarget[0] {
						case '/':
						default:
							symtarget = path + "/" + symtarget
						}
						fi, err = os.Stat(symtarget)
						if err != nil {
							if c.Verbose {
								out("broken symlink target [" + target + "] -> [" + err.Error() + "]")
							}
							fileChan <- &file{symlink: true, syminvalid: true, name: target}
							continue
						}
						if uint32(fi.Mode())&_modeDir != 0 {
							continue // skip symlink dir targets
						}
						if c.Verbose {
							out("symlink [" + target + "] -> [" + symtarget + "]")
						}
						fileChan <- &file{symlink: true, name: target, symtarget: symtarget}
					}
				}
				bg.Done()
			}
		}()
	}
}

// collectTreeWalker ...
func collectTreeWalker() {
	for f := range fileChan {
		switch {
		case f.symlink:
			switch {
			case f.syminvalid:
				syi = append(syi, f.name)
			default:
				sym[f.name] = f.symtarget
			}
		default:
			nfs[f.size] = append(nfs[f.size], f.inode)
			sfs[f.name] = f.size
			ifs[f.name] = f.inode
			total += f.size
			if v, ok := xfs[f.inode]; ok {
				hcount++
				hsave += f.size
				if c.Verbose {
					out("hardlinked [" + f.name + "] -> [" + v[0] + "]")
				}
				xfs[f.inode] = append(xfs[f.inode], f.name)
				continue
			}
			xfs[f.inode] = []string{f.name}
		}
	}
	global.Done()
}
