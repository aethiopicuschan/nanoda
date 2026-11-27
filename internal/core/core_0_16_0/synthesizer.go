package core_0_16_0

import (
	"io"
	stdstrings "strings"
	"unsafe"

	"github.com/aethiopicuschan/nanoda/v2/internal/audioquery"
	"github.com/aethiopicuschan/nanoda/v2/internal/strings"
)

// Synthesizer implements internal/core.Synthesizer for voicevox_core
// 0.16.x. It owns both the native VoicevoxSynthesizer and the OpenJtalkRc
// it was built with (0.16's voicevox_synthesizer_new takes a *borrowed*
// reference to the latter — see voicevox_core.h's own doc comment on
// OpenJtalkRc being reference-counted — so this is the one place that
// reference gets released).
type Synthesizer struct {
	c           *Core
	openJtalk   uintptr
	synthesizer uintptr
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
	if code := s.c.voicevoxVoiceModelFileOpen(vvmPath, &model); code != 0 {
		return s.c.newError("voice_model_file_open: "+vvmPath, code)
	}
	defer s.c.voicevoxVoiceModelFileDelete(model)

	if code := s.c.voicevoxSynthesizerLoadVoiceModel(s.synthesizer, model); code != 0 {
		return s.c.newError("synthesizer_load_voice_model: "+vvmPath, code)
	}
	return nil
}

// ttsOptions mirrors VoicevoxTtsOptions (a single bool, 1 byte) — see
// core.go's package doc comment for why this is passed as a raw uintptr.
type ttsOptions struct {
	enableInterrogativeUpspeak bool
}

// synthesisOptions mirrors VoicevoxSynthesisOptions, which has the same
// single bool as VoicevoxTtsOptions.
type synthesisOptions struct {
	enableInterrogativeUpspeak bool
}

// Tts implements core.Synthesizer. With no params it calls
// voicevox_synthesizer_tts directly; otherwise it goes through an
// AudioQuery, which is what voicevox_synthesizer_tts does internally, so
// both paths produce the same audio for the same AudioQuery.
func (s *Synthesizer) Tts(text string, styleID uint32, params audioquery.Params) (io.ReadCloser, error) {
	if params.IsZero() {
		opts := ttsOptions{enableInterrogativeUpspeak: false}
		var wavLen uintptr
		var wavPtr *byte
		if code := s.c.voicevoxSynthesizerTts(s.synthesizer, text, styleID, *(*uintptr)(unsafe.Pointer(&opts)), &wavLen, &wavPtr); code != 0 {
			return nil, s.c.newError("synthesizer_tts", code)
		}
		return s.newWav(wavPtr, wavLen), nil
	}

	query, err := s.createAudioQuery(text, styleID, params)
	if err != nil {
		return nil, err
	}
	opts := synthesisOptions{enableInterrogativeUpspeak: false}
	var wavLen uintptr
	var wavPtr *byte
	if code := s.c.voicevoxSynthesizerSynthesis(s.synthesizer, query, styleID, *(*uintptr)(unsafe.Pointer(&opts)), &wavLen, &wavPtr); code != 0 {
		return nil, s.c.newError("synthesizer_synthesis", code)
	}
	return s.newWav(wavPtr, wavLen), nil
}

// createAudioQuery returns the AudioQuery JSON with params applied, as a
// copy so that the library's buffer can be freed before returning.
func (s *Synthesizer) createAudioQuery(text string, styleID uint32, params audioquery.Params) (string, error) {
	var queryPtr *byte
	if code := s.c.voicevoxSynthesizerCreateAudioQuery(s.synthesizer, text, styleID, &queryPtr); code != 0 {
		return "", s.c.newError("synthesizer_create_audio_query", code)
	}
	defer s.c.voicevoxJsonFree(uintptr(unsafe.Pointer(queryPtr)))
	return params.Apply(stdstrings.Clone(strings.GoString(queryPtr)))
}

func (s *Synthesizer) newWav(wavPtr *byte, wavLen uintptr) io.ReadCloser {
	raw := unsafe.Slice(wavPtr, wavLen)
	return newWav(raw, func() error {
		s.c.voicevoxWavFree(uintptr(unsafe.Pointer(wavPtr)))
		return nil
	})
}

// Close implements core.Synthesizer.
func (s *Synthesizer) Close() {
	s.c.voicevoxSynthesizerDelete(s.synthesizer)
	s.c.voicevoxOpenJtalkRcDelete(s.openJtalk)
}
