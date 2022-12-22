// package main ...
package main

// import
import (
	"os"

	"paepcke.de/fsdd"
)

// main ...
func main() {
	var err error
	c, l := fsdd.GetDefaultConfig(), len(os.Args)
	switch {
	case l > 1:
		for i := 1; i < l; i++ {
			o := os.Args[i]
			switch o {
			case "--hard-link", "-L":
				c.HardLink = true
				c.Opt += "[--hard-link] "
			case "--sym-link", "-S":
				c.SymLink = true
				c.Opt += "[--sym-link] "
			case "--clean-symlinks", "-C":
				c.ReplaceSymlinks = true
				c.RemoveBrokenSymlinks = true
				c.Opt += "[--clean-symlinks] "
			case "--clean-metadata", "-M":
				c.CleanMetadata = true
				c.Opt += "[--clean-metadata] "
			case "--replace-symlinks", "-R":
				c.ReplaceSymlinks = true
				c.Opt += "[--replace-symlinks] "
			case "--remove-broken-symlinks", "-X":
				c.RemoveBrokenSymlinks = true
				c.Opt += "[--remove-broken-symlinks] "
			case "--fast-hash", "-F":
				c.FastHash = true
				c.Opt += "[--fast-hash] "
			case "--secure-hash", "-H":
				c.FastHash = false
				c.Opt += "[--secure-hash] "
			case "--debug", "-d":
				c.Verbose = true
				c.Debug = true
				c.Opt += "[--debug] "
			case "--verbose", "-v":
				c.Verbose = true
				c.Opt += "[--verbose] "
			case "--help", "-h":
				out(_syntax)
				os.Exit(0)
			default:
				switch {
				case o == ".":
					if c.Path != "" {
						errExit("[more than one target path specified] [" + c.Path + "] and [.]")
					}
					if c.Path, err = os.Getwd(); err != nil {
						errExit("invalid current directory [.] path")
					}
				case isDir(o):
					if c.Path != "" {
						errExit("[more than one target path specified] [" + c.Path + "] and [" + o + "]")
					}
					c.Path = o
				default:
					errExit("[error] [unknown option or invalid path] [" + o + "]")
				}
			}
		}
	default:
		out(_syntax)
		os.Exit(0)
	}
	if c.HardLink && c.SymLink {
		errExit("[hardlink and symlink activated]")
	}
	if c.Path == "" {
		if c.Path, err = os.Getwd(); err != nil {
			errExit("invalid current directory [.] path")
		}
	}
	c.Start()
}

// const
const _syntax string = "syntax: fsdd <start-directory> [options]\n\n--hard-link [-L]\n\t\treplace all duplicated files (within local fs boundary) via hardlinks,\n\t\tsave diskspace and inodes meta data handling, loose duplicate files\n\t\tindividual metadata [not reversible]\n\n--sym-link [-S]\n\t\treplace all duplicated files (within local fs boundary) via symlinks,\n\t\tsave diskspace, keep duplicated individual inodes and metadata\n\n--clean-symlinks [-C]\n\t\treplace all valid resolvable symlinks (local fs boundary) via hardlinks [-R],\n\t\tand remove all broken symlinks [-X], save inodes and metadata handling,\n\t\tloose symlinks metadata [not reversible]\n\n--clean-metadata [-M]\n\t\treset all file and directory timestamps for file create and last-modify date (atime,mtime)\n\t\tto UnixTime 0 (01.01.1970 00:00) and save meta diskspace and handling overhead\n\n--replace-symlinks [-R]\n\t\treplace all valid symlinks via hardlinks, keep broken symlinks\n\t\tsave meta diskspace and handling\n\n--remove-broken-symlinks [-X]\n\t\tdelete all symlinks where the target does not resolve,\n\n--secure-hash [-H]\n\t\tuse SHA512/256 as cryptographic secure hash\n\t\tto avoid intentional designed abusive hash collisions\n\n--fast-hash [-F]\n\t\tfast file content hashing [via MAPHASH instead of SHA512/256]\n\t\tWARNING: fast-hash-deduplication is not 100% intentional [preimage|abuse|collision] resistant!\n\n--verbose [-v]\n--debug [-d]\n--help [-h]\n\nNOTES\n[read-only-mode] Without any activated options (or --verbose only), the application prints only a summary or detailed statistics log."

//
// Little Helper
//

// const
const _modeDir uint32 = 1 << (32 - 1 - 0)

// out ...
func out(message string) {
	os.Stdout.Write([]byte(message + "\n"))
}

// errExit ...
func errExit(message string) {
	out("[error] " + message)
	os.Exit(1)
}

// isDir ...
func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return uint32(fi.Mode())&_modeDir != 0
}
