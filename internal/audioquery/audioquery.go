// Package audioquery overrides fields of the AudioQuery JSON that
// voicevox_synthesizer_create_audio_query returns. Every supported
// voicevox_core version uses the same AudioQuery shape, so this lives
// outside the per-version packages.
package audioquery

import (
	"encoding/json"
	"fmt"
)

// Params holds the AudioQuery fields to override. A nil field keeps the
// value voicevox_core chose; pointers are needed because zero is a
// meaningful value for several of them (pitchScale's default is 0, and an
// intonationScale of 0 means monotone).
type Params struct {
	SpeedScale        *float64
	PitchScale        *float64
	IntonationScale   *float64
	VolumeScale       *float64
	PrePhonemeLength  *float64
	PostPhonemeLength *float64
}

// IsZero reports whether p overrides nothing, in which case callers can
// use voicevox_synthesizer_tts directly instead of the AudioQuery round
// trip.
func (p Params) IsZero() bool {
	return p == Params{}
}

// Apply returns query with p's fields overridden. Fields Params does not
// know about, such as accent_phrases, are carried over unchanged.
func (p Params) Apply(query string) (string, error) {
	if p.IsZero() {
		return query, nil
	}
	// voicevox_core rejects other invalid values with an error, but divides
	// by a zero speedScale and aborts the whole process trying to allocate
	// the resulting length. !(x > 0) also catches NaN.
	if p.SpeedScale != nil && !(*p.SpeedScale > 0) {
		return "", fmt.Errorf("nanoda: speedScale must be positive: %v", *p.SpeedScale)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(query), &fields); err != nil {
		return "", fmt.Errorf("nanoda: parse audio query: %w", err)
	}
	for key, v := range map[string]*float64{
		"speedScale":        p.SpeedScale,
		"pitchScale":        p.PitchScale,
		"intonationScale":   p.IntonationScale,
		"volumeScale":       p.VolumeScale,
		"prePhonemeLength":  p.PrePhonemeLength,
		"postPhonemeLength": p.PostPhonemeLength,
	} {
		if v == nil {
			continue
		}
		raw, err := json.Marshal(*v)
		if err != nil {
			return "", fmt.Errorf("nanoda: %s: %w", key, err)
		}
		fields[key] = raw
	}
	out, err := json.Marshal(fields)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
