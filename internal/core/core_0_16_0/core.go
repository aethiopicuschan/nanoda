// Package core_0_16_0 implements the shape internal/core.Core/Synthesizer
// require against voicevox_core 0.16.x's C API, structurally rather than
// by importing those interfaces directly — doing so would import-cycle
// back to internal/core, which imports this package to dispatch to it.
// internal/core.New adapts *Core's concrete NewSynthesizer return value to
// its own Synthesizer interface at the call site instead. Every native
// call goes through purego (RegisterLibFunc), never cgo — see nanoda's own
// README for why.
//
// # Struct arguments
//
// voicevox_core.h has several functions that take a small options struct
// by value (VoicevoxLoadOnnxruntimeOptions, VoicevoxInitializeOptions,
// VoicevoxTtsOptions). purego only supports passing/returning Go struct
// values on darwin (see purego's func.go: "struct arguments are only
// supported on darwin amd64 & arm64") — everywhere else, including
// Windows and Linux, it panics. Every options struct in this file is
// therefore small enough (<=8 bytes) to fit in a single register, and is
// declared here as a plain Go struct purely to document its layout; the
// actual call always reinterprets it as a uintptr
// (`*(*uintptr)(unsafe.Pointer(&opts))`) before handing it to the
// registered func, which purego treats as an ordinary integer argument.
// This is the same technique nanoda's 0.15 implementation used throughout
// (see the module's git history), not something new to 0.16.
//
// For the same reason, `voicevox_make_default_*_options` functions (which
// return a struct *by value*) can't be called at all on non-darwin
// platforms; this package builds sane defaults in Go instead.
package core_0_16_0

import (
	"fmt"
	"unsafe"

	"github.com/aethiopicuschan/nanoda/v2/internal/strings"
	"github.com/ebitengine/purego"
)

// Core holds the loaded voicevox_core handle and every native function
// pointer purego resolved from it.
type Core struct {
	core uintptr

	voicevoxGetVersion                         func() string
	voicevoxGetOnnxruntimeLibVersionedFilename func() string
	voicevoxOnnxruntimeLoadOnce                func(uintptr, uintptr) int32
	voicevoxOpenJtalkRcNew                     func(string, uintptr) int32
	voicevoxOpenJtalkRcDelete                  func(uintptr)
	voicevoxSynthesizerNew                     func(uintptr, uintptr, uintptr, uintptr) int32
	voicevoxSynthesizerDelete                  func(uintptr)
	voicevoxVoiceModelFileOpen                 func(string, uintptr) int32
	voicevoxVoiceModelFileDelete               func(uintptr)
	voicevoxSynthesizerLoadVoiceModel          func(uintptr, uintptr) int32
	voicevoxSynthesizerTts                     func(uintptr, string, uint32, uintptr, uintptr, uintptr) int32
	voicevoxWavFree                            func(uintptr)
	voicevoxErrorResultToMessage               func(int32) string
}

// New wires up every function this package needs from the already-loaded
// library at the given handle. Registration panics (via purego) if a
// symbol is missing — appropriate here, since internal/core.New only
// reaches this constructor after minimum.Core has already confirmed the
// library reports a 0.16.x version, so a missing symbol means the DLL/so/
// dylib is corrupt or not actually what it claims to be, not a normal
// runtime condition to recover from.
func New(core uintptr) (c *Core) {
	c = &Core{core: core}
	purego.RegisterLibFunc(&c.voicevoxGetVersion, c.core, "voicevox_get_version")
	purego.RegisterLibFunc(&c.voicevoxGetOnnxruntimeLibVersionedFilename, c.core, "voicevox_get_onnxruntime_lib_versioned_filename")
	purego.RegisterLibFunc(&c.voicevoxOnnxruntimeLoadOnce, c.core, "voicevox_onnxruntime_load_once")
	purego.RegisterLibFunc(&c.voicevoxOpenJtalkRcNew, c.core, "voicevox_open_jtalk_rc_new")
	purego.RegisterLibFunc(&c.voicevoxOpenJtalkRcDelete, c.core, "voicevox_open_jtalk_rc_delete")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerNew, c.core, "voicevox_synthesizer_new")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerDelete, c.core, "voicevox_synthesizer_delete")
	purego.RegisterLibFunc(&c.voicevoxVoiceModelFileOpen, c.core, "voicevox_voice_model_file_open")
	purego.RegisterLibFunc(&c.voicevoxVoiceModelFileDelete, c.core, "voicevox_voice_model_file_delete")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerLoadVoiceModel, c.core, "voicevox_synthesizer_load_voice_model")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerTts, c.core, "voicevox_synthesizer_tts")
	purego.RegisterLibFunc(&c.voicevoxWavFree, c.core, "voicevox_wav_free")
	purego.RegisterLibFunc(&c.voicevoxErrorResultToMessage, c.core, "voicevox_error_result_to_message")
	return
}

func (c *Core) GetVersion() string {
	return c.voicevoxGetVersion()
}

// newError turns a VoicevoxResultCode into a Go error carrying the
// library's own message for it (voicevox_error_result_to_message), rather
// than just the bare numeric code.
func (c *Core) newError(op string, code int32) error {
	return fmt.Errorf("nanoda: %s: %s (code %d)", op, c.voicevoxErrorResultToMessage(code), code)
}

// initializeOptions mirrors VoicevoxInitializeOptions
// (acceleration_mode int32 + cpu_num_threads uint16, 8 bytes with
// trailing padding on amd64/arm64) — see the package doc comment for why
// this is passed as a raw uintptr instead of a Go struct value.
type initializeOptions struct {
	accelerationMode int32 // VOICEVOX_ACCELERATION_MODE_AUTO
	cpuNumThreads    uint16
}

// NewSynthesizer implements core.Core. It loads ONNX Runtime (once per
// process; voicevox_onnxruntime_load_once itself is idempotent and
// ignores the options on any call after the first), builds an Open Jtalk
// analyzer from openJtalkDicDir, and constructs the native synthesizer
// with automatic hardware acceleration (mirroring
// voicevox_make_default_initialize_options' own defaults, which can't be
// called directly on non-darwin — see the package doc comment).
func (c *Core) NewSynthesizer(openJtalkDicDir string) (s *Synthesizer, err error) {
	filename := c.voicevoxGetOnnxruntimeLibVersionedFilename()
	// VoicevoxLoadOnnxruntimeOptions is just {filename *const char}: 8
	// bytes, bit-identical to the pointer alone, so no wrapper struct is
	// needed before the uintptr cast.
	var onnxruntime uintptr
	if code := c.voicevoxOnnxruntimeLoadOnce(strings.CString(filename), uintptr(unsafe.Pointer(&onnxruntime))); code != 0 {
		err = c.newError("onnxruntime_load_once", code)
		return
	}

	var openJtalk uintptr
	if code := c.voicevoxOpenJtalkRcNew(openJtalkDicDir, uintptr(unsafe.Pointer(&openJtalk))); code != 0 {
		err = c.newError("open_jtalk_rc_new", code)
		return
	}

	opts := initializeOptions{accelerationMode: 0, cpuNumThreads: 0}
	var synthesizer uintptr
	code := c.voicevoxSynthesizerNew(onnxruntime, openJtalk, *(*uintptr)(unsafe.Pointer(&opts)), uintptr(unsafe.Pointer(&synthesizer)))
	if code != 0 {
		c.voicevoxOpenJtalkRcDelete(openJtalk)
		err = c.newError("synthesizer_new", code)
		return
	}

	s = &Synthesizer{c: c, openJtalk: openJtalk, synthesizer: synthesizer}
	return
}
