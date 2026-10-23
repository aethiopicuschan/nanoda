//go:build darwin || freebsd || linux

package open

import "github.com/ebitengine/purego"

func Open(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}

// HasSymbol reports whether lib exports name. purego.RegisterLibFunc
// panics on a missing symbol, so optional functions added in a patch
// release must be probed with this first.
func HasSymbol(lib uintptr, name string) bool {
	_, err := purego.Dlsym(lib, name)
	return err == nil
}
