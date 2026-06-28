package fsdd

import (
	"os"
	"fmt"
	"strconv"
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

// 
// CONVERTER
//

func itoa(in int) string { return strconv.Itoa(in) }

func itoaU64(in uint64) string { return strconv.FormatUint(in, 10) }

// hruIEC converts value to hru IEC 60027 units
func hruIEC(i uint64, u string) string {
	return hru(i, 1024, u)
}

// hru [human readable units] backend
func hru(i, unit uint64, u string) string {
	if i < unit {
		return fmt.Sprintf("%d %s", i, u)
	}
	div, exp := unit, 0
	for n := i / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	switch u {
	case "":
		return fmt.Sprintf("%.3f %c", float64(i)/float64(div), "kMGTPE"[exp])
	case "bit":
		return fmt.Sprintf("%.0f %c%s", float64(i)/float64(div), "kMGTPE"[exp], u)
	case "bytes", "bytes/sec":
		return fmt.Sprintf("%.1f %c%s", float64(i)/float64(div), "kMGTPE"[exp], u)
	}
	return fmt.Sprintf("%.3f %c%s", float64(i)/float64(div), "kMGTPE"[exp], u)
}
