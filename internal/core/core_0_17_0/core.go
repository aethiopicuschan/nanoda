// Package core_0_17_0 implements the shape internal/core.Core/Synthesizer
// require against voicevox_core 0.17.x's C API, structurally rather than
// by importing those interfaces directly — doing so would import-cycle
// back to internal/core, which imports this package to dispatch to it.
// internal/core.New adapts *Core's concrete NewSynthesizer return value to
// its own Synthesizer interface at the call site instead. Every native
// call goes through purego (RegisterLibFunc), never cgo — see nanoda's own
// README for why.
//
// # Diff against core_0_16_0
//
// 0.17.0's C API is almost identical to 0.16.x's — confirmed against the
// real bundled voicevox_core.h for both versions, function by function, not
// just release-notes prose — except for two breaking changes this package
// exists to absorb:
//
//   - voicevox_get_onnxruntime_lib_versioned_filename was renamed to
//     voicevox_get_onnxruntime_lib_recommended_versioned_filename (the old
//     name is gone from the 0.17.0 header entirely — registering it via
//     purego would panic at startup with a symbol-not-found error, the same
//     failure shape iOS's own onnxruntime split hit before that was fixed).
//     iOS-only (onnxruntime_ios.go, LINK mode) is unaffected — its one
//     function, voicevox_onnxruntime_init_once, kept its exact name.
//   - voicevox_synthesizer_load_voice_model gained a third parameter,
//     struct VoicevoxLoadVoiceModelOptions (a single enum, 4 bytes) — see
//     synthesizer.go's loadVoiceModelOptions for how nanoda passes it.
//
// Everything else this package touches (VoicevoxInitializeOptions,
// VoicevoxTtsOptions, voicevox_synthesizer_tts's own signature, every other
// function name below) is byte-for-byte identical to 0.16.x.
//
// # 0.17.1 additions
//
// 0.17.1 only adds to the C API, so it is handled here rather than in a
// separate package. Its additions are probed with open.HasSymbol and
// registered only when present (see New):
//
//   - voicevox_bytes_free replaces voicevox_wav_free, which still exists
//     but logs a deprecation warning on every call.
//   - voicevox_synthesizer_create_audio_feature/voicevox_synthesizer_render
//     and friends enable streaming synthesis (see stream.go).
//
// # Output parameters
//
// Output parameters are declared as typed pointers (*uintptr, **byte, ...)
// and passed as &local, never as uintptr(unsafe.Pointer(&local)). A
// uintptr hides the pointer from the compiler, so local stays on the
// goroutine stack; if that stack grows inside purego's Go-side call path,
// it is copied elsewhere and the library writes into the stale copy,
// leaving local zero (reliably on a fresh goroutine's first call). Passing
// a typed pointer to a func value makes local escape to the heap instead.
//
// # Struct arguments
//
// voicevox_core.h has several functions that take a small options struct
// by value (VoicevoxInitializeOptions, VoicevoxTtsOptions and the like).
// Each fits in a single register (<=8 bytes), so it is declared here as a
// plain Go struct purely to document its layout and reinterpreted as a
// uintptr (`*(*uintptr)(unsafe.Pointer(&opts))`) at the call, which every
// ABI passes identically to the struct itself. This predates purego's
// struct support (darwin only until v0.10/v0.11) and is kept for
// consistency with core_0_16_0.
//
// VoicevoxAudioFeatureSlice (32 bytes) does not fit in a register, and
// how such a struct is passed differs per ABI (on the stack on SysV
// amd64, by reference on Windows and AAPCS64), so it is passed as a real
// Go struct and purego lays it out per platform.
package core_0_17_0

import (
	"fmt"
	"unsafe"

	"github.com/aethiopicuschan/nanoda/v2/internal/open"
	"github.com/ebitengine/purego"
)

// Core holds the loaded voicevox_core handle and every native function
// pointer purego resolved from it. onnxruntimeFuncs is embedded rather
// than declared here directly — its shape (which functions exist at all)
// differs per platform; see onnxruntime_other.go/onnxruntime_ios.go.
type Core struct {
	core uintptr

	onnxruntimeFuncs

	voicevoxGetVersion                func() string
	voicevoxOpenJtalkRcNew            func(string, *uintptr) int32
	voicevoxOpenJtalkRcDelete         func(uintptr)
	voicevoxSynthesizerNew            func(uintptr, uintptr, uintptr, *uintptr) int32
	voicevoxSynthesizerDelete         func(uintptr)
	voicevoxVoiceModelFileOpen        func(string, *uintptr) int32
	voicevoxVoiceModelFileDelete      func(uintptr)
	voicevoxSynthesizerLoadVoiceModel func(uintptr, uintptr, uintptr) int32
	voicevoxSynthesizerTts            func(uintptr, string, uint32, uintptr, *uintptr, **byte) int32
	voicevoxErrorResultToMessage      func(int32) string

	// bytesFree is voicevox_bytes_free on 0.17.1+ and voicevox_wav_free
	// before it.
	bytesFree func(uintptr)

	streamingFuncs
}

// streamingFuncs are the 0.17.1+ functions streaming synthesis needs.
// They are registered only when the library exports
// voicevox_synthesizer_render; otherwise every field stays nil and
// supportsStreaming is false.
type streamingFuncs struct {
	supportsStreaming bool

	voicevoxSynthesizerCreateAudioQuery   func(uintptr, string, uint32, **byte) int32
	voicevoxStringFree                    func(uintptr)
	voicevoxSynthesizerCreateAudioFeature func(uintptr, string, uint32, uintptr, *uintptr) int32
	voicevoxAudioFeatureLen               func(uintptr) uintptr
	voicevoxAudioFeatureSlice             func(uintptr, uintptr, uintptr, bool, bool, bool, *audioFeatureSlice) int32
	voicevoxSynthesizerRender             func(uintptr, audioFeatureSlice, *uintptr, **byte) int32
	voicevoxAudioFeatureDelete            func(uintptr)
}

// New wires up every function this package needs from the already-loaded
// library at the given handle. Registration panics (via purego) if a
// symbol is missing — appropriate here, since internal/core.New only
// reaches this constructor after minimum.Core has already confirmed the
// library reports a 0.17.x version, so a missing symbol means the DLL/so/
// dylib is corrupt or not actually what it claims to be, not a normal
// runtime condition to recover from.
func New(core uintptr) (c *Core) {
	c = &Core{core: core}
	purego.RegisterLibFunc(&c.voicevoxGetVersion, c.core, "voicevox_get_version")
	c.registerOnnxruntime()
	purego.RegisterLibFunc(&c.voicevoxOpenJtalkRcNew, c.core, "voicevox_open_jtalk_rc_new")
	purego.RegisterLibFunc(&c.voicevoxOpenJtalkRcDelete, c.core, "voicevox_open_jtalk_rc_delete")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerNew, c.core, "voicevox_synthesizer_new")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerDelete, c.core, "voicevox_synthesizer_delete")
	purego.RegisterLibFunc(&c.voicevoxVoiceModelFileOpen, c.core, "voicevox_voice_model_file_open")
	purego.RegisterLibFunc(&c.voicevoxVoiceModelFileDelete, c.core, "voicevox_voice_model_file_delete")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerLoadVoiceModel, c.core, "voicevox_synthesizer_load_voice_model")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerTts, c.core, "voicevox_synthesizer_tts")
	purego.RegisterLibFunc(&c.voicevoxErrorResultToMessage, c.core, "voicevox_error_result_to_message")
	if open.HasSymbol(c.core, "voicevox_bytes_free") {
		purego.RegisterLibFunc(&c.bytesFree, c.core, "voicevox_bytes_free")
	} else {
		purego.RegisterLibFunc(&c.bytesFree, c.core, "voicevox_wav_free")
	}
	if open.HasSymbol(c.core, "voicevox_synthesizer_render") {
		c.registerStreaming()
	}
	return
}

func (c *Core) registerStreaming() {
	purego.RegisterLibFunc(&c.voicevoxSynthesizerCreateAudioQuery, c.core, "voicevox_synthesizer_create_audio_query")
	purego.RegisterLibFunc(&c.voicevoxStringFree, c.core, "voicevox_string_free")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerCreateAudioFeature, c.core, "voicevox_synthesizer_create_audio_feature")
	purego.RegisterLibFunc(&c.voicevoxAudioFeatureLen, c.core, "voicevox_audio_feature_len")
	purego.RegisterLibFunc(&c.voicevoxAudioFeatureSlice, c.core, "voicevox_audio_feature_slice")
	purego.RegisterLibFunc(&c.voicevoxSynthesizerRender, c.core, "voicevox_synthesizer_render")
	purego.RegisterLibFunc(&c.voicevoxAudioFeatureDelete, c.core, "voicevox_audio_feature_delete")
	c.supportsStreaming = true
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
// this is passed as a raw uintptr instead of a Go struct value. Unchanged
// from 0.16.x.
type initializeOptions struct {
	accelerationMode int32 // VOICEVOX_ACCELERATION_MODE_AUTO
	cpuNumThreads    uint16
}

// NewSynthesizer implements core.Core. It loads ONNX Runtime (once per
// process, via whichever of LOAD/LINK mode this platform's voicevox_core
// build actually exports — see loadOnnxruntime's platform-specific
// implementations), builds an Open Jtalk analyzer from openJtalkDicDir,
// and constructs the native synthesizer with automatic hardware
// acceleration (mirroring voicevox_make_default_initialize_options' own
// defaults).
func (c *Core) NewSynthesizer(openJtalkDicDir string) (s *Synthesizer, err error) {
	onnxruntime, err := c.loadOnnxruntime()
	if err != nil {
		return
	}

	var openJtalk uintptr
	if code := c.voicevoxOpenJtalkRcNew(openJtalkDicDir, &openJtalk); code != 0 {
		err = c.newError("open_jtalk_rc_new", code)
		return
	}

	opts := initializeOptions{accelerationMode: 0, cpuNumThreads: 0}
	var synthesizer uintptr
	code := c.voicevoxSynthesizerNew(onnxruntime, openJtalk, *(*uintptr)(unsafe.Pointer(&opts)), &synthesizer)
	if code != 0 {
		c.voicevoxOpenJtalkRcDelete(openJtalk)
		err = c.newError("synthesizer_new", code)
		return
	}

	s = &Synthesizer{c: c, openJtalk: openJtalk, synthesizer: synthesizer}
	return
}
