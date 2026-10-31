//go:build windows

package open

import (
	"golang.org/x/sys/windows"
)

func Open(name string) (uintptr, error) {
	handle, err := windows.LoadDLL(name)
	if err != nil {
		// LoadDLL returns a nil *DLL alongside the error (e.g. when name
		// doesn't exist) — dereferencing handle.Handle below without this
		// check panics instead of letting the caller handle a normal
		// "library not found" error.
		return 0, err
	}
	return uintptr(handle.Handle), nil
}

// HasSymbol reports whether lib exports name. purego.RegisterLibFunc
// panics on a missing symbol, so optional functions added in a patch
// release must be probed with this first.
func HasSymbol(lib uintptr, name string) bool {
	_, err := windows.GetProcAddress(windows.Handle(lib), name)
	return err == nil
}
