//go:build ios

package core_0_16_0

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// onnxruntimeFuncs wires the LINK-mode onnxruntime function
// (VOICEVOX_LINK_ONNXRUNTIME) iOS voicevox_core builds export: ONNX
// Runtime is a normal link-time dependency of voicevox_core itself on
// this platform — resolved by dyld the moment voicevox_core is dlopen'd,
// via its own @rpath/voicevox_onnxruntime.framework/voicevox_onnxruntime
// load command (see example/mobile/ios/README.md's VOICEVOX section for
// how that framework gets embedded/staged). There is no separate library
// for Go code to locate or dlopen — voicevox_onnxruntime_init_once just
// hands back the instance voicevox_core already linked against. See
// onnxruntime_other.go for every other platform's LOAD-mode counterpart,
// and its own doc comment for how this split was confirmed against the
// real bundled header.
type onnxruntimeFuncs struct {
	voicevoxOnnxruntimeInitOnce func(uintptr) int32
}

func (c *Core) registerOnnxruntime() {
	purego.RegisterLibFunc(&c.voicevoxOnnxruntimeInitOnce, c.core, "voicevox_onnxruntime_init_once")
}

// loadOnnxruntime implements the LINK-mode flow: no filename resolution,
// no dlopen — just ask voicevox_core for the instance it already linked
// against (idempotent, like the LOAD-mode side's *_load_once).
func (c *Core) loadOnnxruntime() (uintptr, error) {
	var onnxruntime uintptr
	if code := c.voicevoxOnnxruntimeInitOnce(uintptr(unsafe.Pointer(&onnxruntime))); code != 0 {
		return 0, c.newError("onnxruntime_init_once", code)
	}
	return onnxruntime, nil
}
