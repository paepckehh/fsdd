# fsdd - [f]ile [s]ystem [d]e[d]uplication

[paepche.de/fsdd](https://paepcke.de/fsdd)

- safe space by deduplicate files at (logical) filesystem layer (biggest impact and speedup)
- supports all unix-like-inode-based filesystems
- replace [slow|space|time|intensive] symlinks via [fast|cheap] hardlinks
- perfect for all types of [compressed|read-only|rootfs] filesystems 
- reset all [expensive|noisy|leaking|bad-to-compress] metadata eg, birth-time, last-mod time, last-access time
- analyze exsting filesystem states, backtrack hardlinks, symlinks, savings 
- does NOT cross filesystem boundaries, does not follow any [fs-external] symlinks 
- yes, its fast 
- 100% pure go, minimal exernal imports, use as app or api (see api.go)

EXPLICIT DATALOSS WARNING

- HARDLINKS ARE ABSOLUTE SAFE AND GREAT FOR BUILDING SMALL/FAST READ-ONLY FILESYSTEMS.
- ACTIVATE THE WRITE OPTIONS ON R/W FILESYSTEM ONLY IF YOU ARE 100% SURE YOU UNDERSTAND HARDLINKS!


## SHOWTIME 

### Tame your excessive go mod cache! Even on zfs/btrfs, you will love your new, fast and small go module cache!
``` Shell
cd $GOMODCACHE && fsdd --hard-link . 
FSDD [start] [/usr/store/go]  [hash:xxh3] 
FSDD [_done] [time: 41.200862ms]
FSDD [stats] [files:13329] [inode(s): 8680] [sym.valid: 0] [sym.invalid: 0] [data blocks: 277.3 Mbytes]
FSDD [_info] [new deduplication savings] [inode(s): 4649] [data blocks: 61.7 Mbytes]
``` 

### Same with detailed file listing.
``` Shell
fsdd --verbose . 
[...]
```

### Replace all duplicates.

``` Shell
fsdd --hard-link . 
[...]
```

### More?

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
		use Blake3 as cryptographic secure content hash function
		to avoid intentional designed abusive hash collisions

 --fast-hash [-F]
		10x faster and lightweight file content hashing [via xxh3/64 instead of Blake3]
		WARNING: fast-hash-deduplication is not 100% intentional [preimage|abuse|collision] resistant
```
