package core_0_17_0

import (
	"io"
	"unsafe"
)

// Synthesizer implements internal/core.Synthesizer for voicevox_core
// 0.17.x. It owns both the native VoicevoxSynthesizer and the OpenJtalkRc
// it was built with (voicevox_synthesizer_new takes a *borrowed* reference
// to the latter — see voicevox_core.h's own doc comment on OpenJtalkRc
// being reference-counted — so this is the one place that reference gets
// released).
type Synthesizer struct {
	c           *Core
	openJtalk   uintptr
	synthesizer uintptr
}

// loadVoiceModelOptions mirrors VoicevoxLoadVoiceModelOptions — new in
// 0.17.0, a single VoicevoxOnExistingVoiceModelId enum (4 bytes; see the
// package doc comment for why this is passed as a raw uintptr instead of a
// Go struct value). onExisting: 0 is
// VOICEVOX_ON_EXISTING_VOICE_MODEL_ID_ERROR, the documented default and
// exactly what 0.16.x's load_voice_model (no options parameter at all)
// always did — loading the same model twice errors instead of silently
// reloading it. Kept at the zero value deliberately: nanoda doesn't expose
// RELOAD/SKIP yet, so every caller gets 0.16-equivalent behavior.
type loadVoiceModelOptions struct {
	onExisting int32
}

// LoadVoiceModel implements core.Synthesizer. It opens the .vvm file,
// loads it into the synthesizer, then immediately closes the file handle
// — VoicevoxVoiceModelFile represents the *file*, not the model data the
// synthesizer now holds, and the official Python API's own usage pattern
// (VoiceModelFile.open() used as a context manager around exactly one
// load_voice_model call) confirms this is the intended lifetime rather
// than something that needs to outlive the load.
func (s *Synthesizer) LoadVoiceModel(vvmPath string) error {
	var model uintptr
	if code := s.c.voicevoxVoiceModelFileOpen(vvmPath, uintptr(unsafe.Pointer(&model))); code != 0 {
		return s.c.newError("voice_model_file_open: "+vvmPath, code)
	}
	defer s.c.voicevoxVoiceModelFileDelete(model)

	opts := loadVoiceModelOptions{onExisting: 0}
	if code := s.c.voicevoxSynthesizerLoadVoiceModel(s.synthesizer, model, *(*uintptr)(unsafe.Pointer(&opts))); code != 0 {
		return s.c.newError("synthesizer_load_voice_model: "+vvmPath, code)
	}
	return nil
}

// ttsOptions mirrors VoicevoxTtsOptions (a single bool, 1 byte) — see
// core.go's package doc comment for why this is passed as a raw uintptr.
// Unchanged from 0.16.x.
type ttsOptions struct {
	enableInterrogativeUpspeak bool
}

// Tts implements core.Synthesizer.
func (s *Synthesizer) Tts(text string, styleID uint32) (io.ReadCloser, error) {
	opts := ttsOptions{enableInterrogativeUpspeak: false}
	var wavLen uintptr
	// wavPtr is deliberately a real *byte, not a uintptr: that lets the
	// slice below read it directly (unsafe.Slice(wavPtr, wavLen)) instead
	// of round-tripping through unsafe.Pointer(uintptr), which go vet
	// flags as a possible misuse (it can't verify a uintptr obtained from
	// C still denotes a live pointer). Passing &wavPtr to the C call
	// converts pointer->uintptr, the direction vet doesn't flag.
	var wavPtr *byte
	code := s.c.voicevoxSynthesizerTts(s.synthesizer, text, styleID, *(*uintptr)(unsafe.Pointer(&opts)), uintptr(unsafe.Pointer(&wavLen)), uintptr(unsafe.Pointer(&wavPtr)))
	if code != 0 {
		return nil, s.c.newError("synthesizer_tts", code)
	}
	raw := unsafe.Slice(wavPtr, wavLen)
	return newWav(raw, func() error {
		s.c.bytesFree(uintptr(unsafe.Pointer(wavPtr)))
		return nil
	}), nil
}

// Close implements core.Synthesizer.
func (s *Synthesizer) Close() {
	s.c.voicevoxSynthesizerDelete(s.synthesizer)
	s.c.voicevoxOpenJtalkRcDelete(s.openJtalk)
}
