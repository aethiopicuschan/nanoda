//go:build darwin || freebsd || linux

package open

import "github.com/ebitengine/purego"

func Open(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}
