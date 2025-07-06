package fsdd

import (
	"runtime"
	"sync"
	"time"
)

type file struct {
	name          string   // full qualified file name (path/name)
	size          uint64   // data size in bytes
	inode         uint64   // inode number
	hash          [32]byte // file hash
	fasthash      uint64   // file fast hash
	symlink       bool     // true if file is a symlink
	syminvalid    bool     // true if file is a symlink & invalid
	newSymLink    bool     // true if new symlink is requested
	symtarget     string   // pre-existing symlink target
	newlinktarget string   // target filename (hard|sym)
}

var (
	// global config
	c *Config
	// fs state analysis
	xfs = make(map[uint64][]string)   // map inode to []name(s)
	nfs = make(map[uint64][]uint64)   // map size to inode
	sfs = make(map[string]uint64)     // map name to size
	ifs = make(map[string]uint64)     // map name to inode
	hfs = make(map[[32]byte][]string) // map hash to name
	hFs = make(map[uint64][]string)   // map fasthash to name
	sym = make(map[string]string)     // map existing symlinks
	syi []string                      // array of invalid symlinks
	// global channel
	fileChan  = make(chan *file, 10000)
	hashChan  = make(chan *file, 10000)
	linkChan  = make(chan *file, 10000)
	failChan  = make(chan *file, 10000)
	metaChan  = make(chan string, 10000)
	dirChan   = make(chan string, 10000)
	namesChan = make(chan string, 10000)
	rmsymChan = make(chan string, 10000)
	// global locks
	bg, ctl, global sync.WaitGroup
	// global counter
	total, hsave, hcount, stotal, sinodes, ssym, failsize, failinodes, failssym uint64
)

func (config *Config) run() {
	// init options (r/o) via global
	c = optionsSanityCheck(config)

	// tick
	t0 := time.Now()
	reportInit()

	// config
	worker := runtime.NumCPU()

	// recursive tree walk
	{
		// spin up collector
		global.Add(1)
		go collectTreeWalker()

		// spin up recursive file tree walker
		if c.CleanMetadata {
			ctl.Add(worker)
			go cleanMetadata(worker)
		}
		go workerTreeWalker(worker)

		// init treeWalk
		bg.Add(1)
		dirChan <- c.Path

		// wait for recursive treewalk
		bg.Wait()
		close(dirChan)
		close(fileChan)
		close(metaChan)
		ctl.Wait()
	}
	global.Wait()

	// file hashing
	{
		// spin up hashsum worker collector
		global.Add(1)
		go collectFileHash()

		// spin up hashsum worker
		bg.Add(worker)
		go workerFileHash(worker)

		// feed same-size-files into hashsum worker
		go feedFileHash()

		// wait
		bg.Wait()
		close(hashChan)
	}
	global.Wait()

	// apply changes
	{
		// collect stats
		global.Add(1)
		go collectStats()

		// spin up file system change worker
		switch {
		case c.HardLink || c.SymLink || c.ReplaceSymlinks:
			bg.Add(worker)
			go workerLink(worker)
		case c.RemoveBrokenSymlinks:
			bg.Add(worker)
			go workerRemoveBrokenSymlinks(worker)
		}

		// feed identified dupes & clean targets into linker/remove/stats channel
		ctl.Add(1)
		go feedLinkStats()
		if c.ReplaceSymlinks {
			ctl.Add(1)
			go feederReplaceSymlinks()
		}
		if c.RemoveBrokenSymlinks {
			ctl.Add(1)
			go feederRemoveBrokenSymlinks()
		}

		// wait
		ctl.Wait()
		close(rmsymChan)
		close(linkChan)
		bg.Wait()
		close(failChan)
	}
	global.Wait()

	// report summary when finished
	reportSummary(t0)
}

//
// LITTLE HELPER
//

// optionsSanityCheck ...
func optionsSanityCheck(config *Config) *Config {
	config.Path = verifyPath(pathSanitizer(config.Path))
	if config.HardLink && config.SymLink {
		errExit("[hardlink and symlink activated] [please choose only one target method]")
	}
	return config
}
