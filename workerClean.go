package fsdd

import (
	"os"
	"syscall"
)

func cleanMetadata(worker int) {
	for i := 0; i < worker; i++ {
		var utimes [2]syscall.Timespec
		utimes[0] = syscall.NsecToTimespec(0)
		utimes[1] = utimes[0]
		go func() {
			for target := range metaChan {
				if err := syscall.UtimesNano(target, utimes[:]); err != nil {
					errOut("[unable to clean meta] [atime,mtime] [" + target + "] [" + err.Error() + "]")
				}
			}
			ctl.Done()
		}()
	}
}

func workerRemoveBrokenSymlinks(worker int) {
	for i := 0; i < worker; i++ {
		go func() {
			for s := range rmsymChan {
				err := os.Remove(s)
				if err != nil {
					errOut("[unable to remove invalid symlink] [" + s + "] [" + err.Error() + "]")
					failChan <- &file{syminvalid: true}
					continue
				}
				if c.Verbose {
					out("[broken symlink removed] [" + s + "]")
				}
			}
			bg.Done()
		}()
	}
}

func feederReplaceSymlinks() {
	for k, v := range sym {
		linkChan <- &file{name: v, newlinktarget: k}
		ssym++
	}
	ctl.Done()
}

func feederRemoveBrokenSymlinks() {
	for _, s := range syi {
		rmsymChan <- s
	}
	ctl.Done()
}
