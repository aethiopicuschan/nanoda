package core

import (
	"io"
	"strings"

	"github.com/aethiopicuschan/nanoda/v2/internal/core/core_0_16_0"
	"github.com/aethiopicuschan/nanoda/v2/internal/core/minimum"
)

// Core is the version-agnostic surface every internal/core/core_x_y_z
// implementation must provide. New dispatches to the right one based on
// the loaded library's own reported version (minimum.GetVersion), so
// callers elsewhere in nanoda never need to know which voicevox_core
// version they're actually talking to.
type Core interface {
	GetVersion() string
	// NewSynthesizer builds everything needed to synthesize speech from
	// openJtalkDicDir (an Open JTalk system dictionary directory): loading
	// the version-appropriate ONNX Runtime if that version needs one,
	// constructing the Open JTalk analyzer, and constructing the native
	// synthesizer object.
	NewSynthesizer(openJtalkDicDir string) (Synthesizer, error)
}

// Synthesizer is the version-agnostic synthesis surface. Model loading and
// TTS are the only operations current support needs; AudioQuery/
// AccentPhrase/UserDict are deliberately not part of this interface yet —
// see nanoda's README "対応状況" for what's still missing versus the 0.15
// implementation.
type Synthesizer interface {
	// LoadVoiceModel loads a single .vvm file into the synthesizer.
	LoadVoiceModel(vvmPath string) error
	// Tts synthesizes text as the given style, returning a WAV stream the
	// caller must Close.
	Tts(text string, styleID uint32) (io.ReadCloser, error)
	// Close releases every native resource this Synthesizer owns.
	Close()
}

// core016Adapter adapts *core_0_16_0.Core to the Core interface.
// core_0_16_0.Core.NewSynthesizer returns its own concrete *Synthesizer
// type rather than this package's Synthesizer interface (see
// core_0_16_0's own package doc comment for why: importing this package
// from there to name the interface directly would create an import
// cycle). A concrete return type doesn't structurally satisfy an
// interface method that returns an interface type — Go requires an exact
// signature match — so this thin wrapper performs the interface
// conversion at the one call site that's allowed to know about both
// sides.
type core016Adapter struct {
	*core_0_16_0.Core
}

func (a core016Adapter) NewSynthesizer(openJtalkDicDir string) (Synthesizer, error) {
	return a.Core.NewSynthesizer(openJtalkDicDir)
}

// New picks the internal/core/core_x_y_z implementation matching the
// library's own reported version, or returns nil if none matches.
//
// The match is against the version's "X.Y" prefix, not an exact string:
// voicevox_get_version() returns the full release (e.g. "0.16.4"), and a
// single core_0_16_0 implementation is expected to handle every 0.16.x
// patch release, not just one exact one.
func New(core uintptr) (c Core) {
	mc := minimum.New(core)
	switch {
	case strings.HasPrefix(mc.GetVersion(), "0.16."):
		c = core016Adapter{core_0_16_0.New(core)}
		return
	}
	return
}
