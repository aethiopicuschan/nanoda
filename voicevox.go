package nanoda

import (
	"io"
	"os"
	"path/filepath"

	"github.com/aethiopicuschan/nanoda/v2/internal/core"
	"github.com/aethiopicuschan/nanoda/v2/internal/open"
)

// StyleId identifies a character's speech style, matching
// VoicevoxStyleId (a plain uint32) in voicevox_core.h.
type StyleId uint32

type Voicevox struct {
	core core.Core
	// openJtalkPath/modelPath are kept rather than acted on immediately,
	// so NewSynthesizer (built lazily, on demand) and LoadAllModels can
	// each use them without New itself needing to construct a synthesizer
	// nobody may end up wanting.
	openJtalkPath string
	modelPath     string
}

func New(corePath string, openJtalkPath string, modelPath string) (v *Voicevox, err error) {
	c, err := open.Open(corePath)
	if err != nil {
		return
	}
	core := core.New(c)
	if core == nil {
		err = ErrUnsupportedVersion
		return
	}
	v = &Voicevox{
		core:          core,
		openJtalkPath: openJtalkPath,
		modelPath:     modelPath,
	}
	return
}

func (v *Voicevox) Version() string {
	return v.core.GetVersion()
}

// Synthesizer performs text-to-speech. Build one with
// (*Voicevox).NewSynthesizer.
type Synthesizer struct {
	inner     core.Synthesizer
	modelPath string
}

// NewSynthesizer builds a Synthesizer using the Open JTalk dictionary
// directory given to New. Each call constructs its own native synthesizer
// (and, on 0.16.x, its own Open JTalk analyzer) — see (*Synthesizer).Close.
func (v *Voicevox) NewSynthesizer() (s *Synthesizer, err error) {
	inner, err := v.core.NewSynthesizer(v.openJtalkPath)
	if err != nil {
		return
	}
	s = &Synthesizer{inner: inner, modelPath: v.modelPath}
	return
}

// LoadVoiceModel loads a single .vvm file.
func (s *Synthesizer) LoadVoiceModel(vvmPath string) error {
	return s.inner.LoadVoiceModel(vvmPath)
}

// LoadAllModels loads every .vvm file directly under the modelPath given
// to New.
func (s *Synthesizer) LoadAllModels() (err error) {
	entries, err := os.ReadDir(s.modelPath)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".vvm" {
			continue
		}
		if err = s.LoadVoiceModel(filepath.Join(s.modelPath, entry.Name())); err != nil {
			return
		}
	}
	return
}

// Tts synthesizes text as the given style, returning a WAV stream the
// caller must Close.
func (s *Synthesizer) Tts(text string, styleID StyleId) (io.ReadCloser, error) {
	return s.inner.Tts(text, uint32(styleID))
}

// Close releases every native resource this Synthesizer owns.
func (s *Synthesizer) Close() {
	s.inner.Close()
}
