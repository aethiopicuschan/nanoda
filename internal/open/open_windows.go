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
