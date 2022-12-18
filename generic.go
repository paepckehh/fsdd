// package fsdd ...
package fsdd

// import
import (
	"fmt"
	"strconv"
)

// itoa ...
func itoa(in int) string { return strconv.Itoa(in) }

// itoaU64 ...
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
