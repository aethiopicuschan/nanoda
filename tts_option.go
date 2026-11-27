package nanoda

import "github.com/aethiopicuschan/nanoda/v2/internal/audioquery"

// TtsOption adjusts the AudioQuery that Tts and TtsStream synthesize from.
// An option that is not given keeps the value voicevox_core chose for the
// text and style. Values other than a non-positive speed are passed to
// voicevox_core as is: it returns an error for a negative volume or
// silence length, and an extreme value can make it fail to allocate memory
// for the audio.
type TtsOption func(*audioquery.Params)

// WithSpeedScale sets the speaking speed. 1 is normal; larger is faster.
// It must be positive.
func WithSpeedScale(v float64) TtsOption {
	return func(p *audioquery.Params) { p.SpeedScale = &v }
}

// WithPitchScale shifts the pitch. 0 is normal; positive is higher. The
// VOICEVOX editor offers -0.15 to 0.15.
func WithPitchScale(v float64) TtsOption {
	return func(p *audioquery.Params) { p.PitchScale = &v }
}

// WithIntonationScale scales the intonation. 1 is normal, 0 is monotone,
// and larger is more expressive.
func WithIntonationScale(v float64) TtsOption {
	return func(p *audioquery.Params) { p.IntonationScale = &v }
}

// WithVolumeScale scales the volume. 1 is normal.
func WithVolumeScale(v float64) TtsOption {
	return func(p *audioquery.Params) { p.VolumeScale = &v }
}

// WithPrePhonemeLength sets the length of the silence before the speech,
// in seconds.
func WithPrePhonemeLength(v float64) TtsOption {
	return func(p *audioquery.Params) { p.PrePhonemeLength = &v }
}

// WithPostPhonemeLength sets the length of the silence after the speech,
// in seconds.
func WithPostPhonemeLength(v float64) TtsOption {
	return func(p *audioquery.Params) { p.PostPhonemeLength = &v }
}

func newParams(opts []TtsOption) (p audioquery.Params) {
	for _, opt := range opts {
		opt(&p)
	}
	return
}
