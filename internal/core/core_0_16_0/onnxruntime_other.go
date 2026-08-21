//go:build !ios

package core_0_16_0

import (
	"unsafe"

	"github.com/aethiopicuschan/nanoda/v2/internal/strings"
	"github.com/ebitengine/purego"
)

// onnxruntimeFuncs wires the LOAD-mode onnxruntime functions
// (VOICEVOX_LOAD_ONNXRUNTIME) every non-iOS voicevox_core build exports:
// the ONNX Runtime shared library is a separate file the caller must
// locate and dlopen explicitly. See onnxruntime_ios.go for iOS's
// LINK-mode counterpart — voicevox_core.h documents these two modes as
// mutually exclusive per platform (iOS ships only
// VOICEVOX_LINK_ONNXRUNTIME; every other platform ships only
// VOICEVOX_LOAD_ONNXRUNTIME), confirmed against the real header bundled
// in voicevox_core-ios-xcframework-cpu-0.16.4's own Headers/voicevox_core.h.
type onnxruntimeFuncs struct {
	voicevoxGetOnnxruntimeLibVersionedFilename func() string
	voicevoxOnnxruntimeLoadOnce                func(uintptr, uintptr) int32
}

func (c *Core) registerOnnxruntime() {
	purego.RegisterLibFunc(&c.voicevoxGetOnnxruntimeLibVersionedFilename, c.core, "voicevox_get_onnxruntime_lib_versioned_filename")
	purego.RegisterLibFunc(&c.voicevoxOnnxruntimeLoadOnce, c.core, "voicevox_onnxruntime_load_once")
}

// loadOnnxruntime implements the LOAD-mode flow: ask the library for its
// own expected ONNX Runtime filename, then dlopen it by that name (once
// per process — voicevox_onnxruntime_load_once is idempotent and ignores
// the options on any call after the first).
func (c *Core) loadOnnxruntime() (uintptr, error) {
	filename := c.voicevoxGetOnnxruntimeLibVersionedFilename()
	// VoicevoxLoadOnnxruntimeOptions is just {filename *const char}: 8
	// bytes, bit-identical to the pointer alone, so no wrapper struct is
	// needed before the uintptr cast (see core.go's package doc comment on
	// why options structs are passed this way).
	var onnxruntime uintptr
	if code := c.voicevoxOnnxruntimeLoadOnce(strings.CString(filename), uintptr(unsafe.Pointer(&onnxruntime))); code != 0 {
		return 0, c.newError("onnxruntime_load_once", code)
	}
	return onnxruntime, nil
}
