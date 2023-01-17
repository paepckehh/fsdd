# OVERVIEW 
[![Go Reference](https://pkg.go.dev/badge/paepcke.de/dnsresolver.svg)](https://pkg.go.dev/paepcke.de/dnsresolver) [![Go Report Card](https://goreportcard.com/badge/paepcke.de/fsdd)](https://goreportcard.com/report/paepcke.de/fsdd) [![Go Build](https://github.com/paepckehh/fsdd/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/fsdd/actions/workflows/golang.yml)

[paepche.de/fsdd](https://paepcke.de/fsdd/)

fsdd - [f]ile [s]ystem [d]e[d]uplication

- deduplicate files at (logical) filesystem layer (simple/fast)
- supports unix-like-inode-based filesystems (no support for microsoft windows yet))
- replace [slow|space|time|intensive] symlinks via [fast|cheap] hardlinks
- perfect for all types of [compressed|read-only|rootfs|cache] filesystems 
- reset all [expensive|noisy|leaking|bad-to-compress] metadata eg, birth-time, last-mod time, last-access time
- analyze exsting filesystem states, backtrack hardlinks, symlinks, savings 
- does NOT cross filesystem boundaries, does not follow any [fs-external] symlinks 
- yes, its fast (uses the new maphash runtime package)
- 100% pure go, stdlib only, dependency free, use as app or api (see api.go)


# INSTALL

```
go install paepcke.de/fsdd/cmd/fsdd@latest
```

### DOWNLOAD (prebuild)

[github.com/paepckehh/fsdd/releases](https://github.com/paepckehh/fsdd/releases)

# SHOWTIME 

## Tame your excessive go mod cache! Even on zfs/btrfs, you will love your new, fast and small go module cache!
``` Shell
cd $GOMODCACHE && fsdd --hard-link . 
FSDD [start] [/usr/store/go]  [hash:MAPHASH] 
FSDD [_done] [time: 39.000221ms]
FSDD [stats] [files:13329] [inode(s): 8680] [sym.valid: 0] [sym.invalid: 0] [data blocks: 277.1 Mbytes]
FSDD [_info] [new deduplication savings] [inode(s): 4649] [data blocks: 66.9 Mbytes]
``` 

## Same with detailed file listing.
``` Shell
fsdd --verbose . 
[...]
```

## Replace all duplicates.

``` Shell
fsdd --hard-link . 
[...]
```

## More?

``` Shell
fsdd --help 
syntax: fsdd <start-directory> [options]

 --hard-link [-L]
		replace all duplicated files (within local fs boundary) via hardlinks,
		save diskspace and inodes meta data handling, loose duplicate files
		individual metadata [not reversible]

 --sym-link [-S]
		replace all duplicated files (within local fs boundary) via symlinks,
		save diskspace, keep duplicated individual inodes and metadata

 --clean-symlinks [-C]
		replace all valid resolvable symlinks (local fs boundary) via hardlinks [-R],
		and remove all broken symlinks [-X], save inodes and metadata handling,
 		loose symlinks metadata [not reversible]

 --clean-metadata [-M]
		reset all file and directory timestamps for file create and last-modify date (atime,mtime)
		to UnixTime 0 (01.01.1970 00:00) and save meta diskspace and handling overhead

 --replace-symlinks [-R]
		replace all valid symlinks via hardlinks, keep broken symlinks
		save meta diskspace and handling

 --remove-broken-symlinks [-X]
		delete all symlinks where the target does not resolve,

 --secure-hash [-H]
		use SHA512/256 as cryptographic secure hash
		to avoid intentional designed abusive hash collisions

 --fast-hash [-F]
		extreme fast file content hashing [via MAPHASH instead of SHA512/256]
		WARNING: fast-hash-deduplication is not 100% intentional [preimage|abuse|collision] resistant!
```

# WARNING

Hardlinks are absolute great for building fast, small and snappy read-only filesystems and most 
types of filesytem backed (automatically) managed caches. But on normal manually managed full read
and write filesystems hardlinks could result in unexpected data changes or data loss.

# DOCS

[pkg.go.dev/paepcke.de/fsdd](https://pkg.go.dev/paepcke.de/fsdd)

# CONTRIBUTION

Yes, Please! PRs Welcome! 
