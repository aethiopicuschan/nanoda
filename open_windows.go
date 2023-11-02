//go:build windows

package nanoda

import (
	"golang.org/x/sys/windows"
)

func openLibrary(name string) (uintptr, error) {
	handle, err := windows.LoadDLL(name)
	return uintptr(handle.Handle), err
}
