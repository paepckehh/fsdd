package fsdd

import (
	"crypto/sha512"
	"hash/maphash"
	"io"
	"os"
)

//
// WORKER
//

func workerFileHash(worker int) {
	switch c.FastHash {
	case true:
		for i := 0; i < worker; i++ {
			go func() {
				for target := range namesChan {
					hashChan <- &file{name: target, fasthash: fastHash(target)}
				}
				bg.Done()
			}()
		}
	case false:
		for i := 0; i < worker; i++ {
			go func() {
				for target := range namesChan {
					hashChan <- &file{name: target, hash: hash(target)}
				}
				bg.Done()
			}()
		}
	}
}

func feedFileHash() {
	dupes := make(map[uint64]bool) // avoid hashing inode dupes
	for _, v := range nfs {
		l := len(v)
		if l > 1 {
			for i := 0; i < l; i++ {
				node := v[i]
				if _, ok := dupes[node]; ok {
					continue
				}
				dupes[node] = true
				namesChan <- xfs[v[i]][0]
			}
		}
	}
	close(namesChan)
}

func collectFileHash() {
	switch c.FastHash {
	case true:
		go collectorFastHash()
	default:
		go collectorHash()
	}
}

func collectorFastHash() {
	for f := range hashChan {
		hFs[f.fasthash] = append(hFs[f.fasthash], f.name)
	}
	global.Done()
}

func collectorHash() {
	for f := range hashChan {
		hfs[f.hash] = append(hfs[f.hash], f.name)
	}
	global.Done()
}

//
// BACKEND
//

const (
	_hashSize      = 32
	_hashBlockSize = 1024 * 32
)

// hash a file via sha512/256
// fast enough on most modern 64bit arm64/x86-64 systems with sha assisted hardware instruction set
func hash(file string) [_hashSize]byte {
	f, _ := os.Open(file) // access already verified, skip double check here
	r, h := io.Reader(f), sha512.New512_256()
	for {
		block := make([]byte, _hashBlockSize)
		l, _ := r.Read(block)
		if l < _hashBlockSize {
			_, err := h.Write(block)
			if err != nil {
				panic("[internal error] [unable to continue] [hash] [state]")
			}
			break
		}
		_, err := h.Write(block)
		if err != nil {
			panic("[internal error] [unable to continue] [hash] [state]")
		}
	}
	f.Close()
	var hashOut [_hashSize]byte
	for k, v := range h.Sum(nil) {
		hashOut[k] = v
	}
	return hashOut
}

// mseed intit an static seed
var mseed = maphash.MakeSeed()

// fasthash hash a via the new maphash pkg
// extreme fast, but not secure against intentional crafted collisions (exact filesize & filehash must meet here, hard to archive)
func fastHash(file string) uint64 {
	f, _ := os.Open(file) // access already verified, skip double check here
	r := io.Reader(f)
	var h maphash.Hash
	h.SetSeed(mseed)
	for {
		block := make([]byte, _hashBlockSize)
		l, _ := r.Read(block)
		if l < _hashBlockSize {
			_, err := h.Write(block)
			if err != nil {
				panic("[internal error] [unable to continue] [hash] [state] [maphash]")
			}
			break
		}
		_, err := h.Write(block)
		if err != nil {
			panic("[internal error] [unable to continue] [hash] [state] [maphash]")
		}
	}
	f.Close()
	return h.Sum64()
}
