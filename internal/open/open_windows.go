//go:build windows

package open

import (
	"golang.org/x/sys/windows"
)

func Open(name string) (uintptr, error) {
	handle, err := windows.LoadDLL(name)
	return uintptr(handle.Handle), err
}
