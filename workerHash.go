// package fsdd
package fsdd

// import
import (
	"io"
	"os"
	"crypto/sha512"

	// "github.com/zeebo/blake3"
	"github.com/zeebo/xxh3"
)

//
// WORKER
//

// workerFileHash ...
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

// feedFileHash ...
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

// collectFileHash ...
func collectFileHash() {
	switch c.FastHash {
	case true:
		go collectorFastHash()
	default:
		go collectorHash()
	}
}

// collectorFastHash ...
func collectorFastHash() {
	for f := range hashChan {
		hFs[f.fasthash] = append(hFs[f.fasthash], f.name)
	}
	global.Done()
}

// collectorHash ...
func collectorHash() {
	for f := range hashChan {
		hfs[f.hash] = append(hfs[f.hash], f.name)
	}
	global.Done()
}

//
// BACKEND
//

// const
const (
	_hashSize      = 32
	_hashBlockSize = 1024 * 32
)

// hash hashes a file  [crypto|preimage|collision] resistant & secure -> sha512 [fast enough for almost any storage/cpu combo]
func hash(file string) [_hashSize]byte {
	f, _ := os.Open(file)
	r, h := io.Reader(f), sha512.New()
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
		h.Write(block)
	}
	f.Close()
	var hashOut [_hashSize]byte
	for k, v := range h.Sum(nil) {
		hashOut[k] = v
	}
	if c.Debug {
		out("[debug] [hashing] [sha512] [" + file + "]")
	}
	return hashOut
}

// fasthash hahes a file [collision] secure -> xxH3/128 [ultra fast/light, not secure against *intentional* crafted collisions!]
func fastHash(file string) uint64 {
	f, _ := os.Open(file)
	r, h := io.Reader(f), xxh3.New()
	for {
		block := make([]byte, _hashBlockSize)
		l, _ := r.Read(block)
		if l < _hashBlockSize {
			_, err := h.Write(block)
			if err != nil {
				panic("[internal error] [unable to continue] [hash] [state] [xxh3/64]")
			}
			break
		}
		h.Write(block)
	}
	f.Close()
	if c.Debug {
		out("[debug] [hashing] [xxh3] [" + file + "]")
	}
	return h.Sum64()
}
