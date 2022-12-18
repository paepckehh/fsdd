// package fsdd ...
package fsdd

// import
import (
	"os"
)

//
// Display IO
//

// out ...
func out(message string) {
	os.Stdout.Write([]byte(message + "\n"))
}

//
// Error Display IO
//

// errOut ...
func errOut(m string) {
	out("[error] " + m)
}

// errExit ...
func errExit(m string) {
	errOut(m)
	os.Exit(1)
}

//
// File IO
//

// pathSanitizer ...
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

// verifyPath ...
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
