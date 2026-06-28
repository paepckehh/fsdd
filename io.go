package fsdd

import (
	"fmt"
	"math/bits"
	"os"
)

//
// DISPLAY IO
//

func out(message string) {
	os.Stdout.Write([]byte(message + "\n"))
}

//
// ERROR DISPLAY IO
//

func errOut(m string) {
	out("[error] " + m)
}

func errExit(m string) {
	errOut(m)
	os.Exit(1)
}

// humanUint64, format uint64 into human readable numbers
func humanUint64(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d entries", bytes)
	}
	base := uint(bits.Len64(bytes) / 10)
	val := float64(bytes) / float64(uint64(1<<(base*10)))
	return fmt.Sprintf("%.1f %ci entries", val, " KMGTPE"[base])
}

//
// FILE IO
//

func pathSanitizer(path string) string {
	var err error
	if path == "." {
		path, err = os.Getwd()
		if err != nil {
			errExit("[pathSanitizer] [invalid current dir]")
		}
	}
	return path
}

func verifyPath(path string) string {
	var err error
	_, err = os.ReadDir(path)
	if err != nil {
		errExit("[read root dir] [" + path + "] [" + err.Error() + "]")
	}
	_, err = os.Stat(path)
	if err != nil {
		errExit("[stat root dir] [" + path + "] [" + err.Error() + "]")
	}
	return path
}
