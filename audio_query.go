package nanoda

import (
	"encoding/json"
	"unsafe"

	"github.com/aethiopicuschan/nanoda/internal/strings"
)

// 音声合成用のクエリ
type AudioQuery struct {
	AccentPhrases      []AccentPhrase `json:"accent_phrases"`
	SpeedScale         float64        `json:"speed_scale"`
	PitchScale         float64        `json:"pitch_scale"`
	IntonationScale    float64        `json:"intonation_scale"`
	VolumeScale        float64        `json:"volume_scale"`
	PrePhonemeLength   float64        `json:"pre_phoneme_length"`
	PostPhonemeLength  float64        `json:"post_phoneme_length"`
	OutputSamplingRate int            `json:"output_sampling_rate"`
	OutputStereo       bool           `json:"output_stereo"`
	Kana               string         `json:"kana"`
}

// 音声合成用のクエリを作成する
func (s *Synthesizer) createAudioQuery(text string, styleID StyleId, enableKana bool) (a AudioQuery, err error) {
	var ptr *byte
	var code ResultCode
	if enableKana {
		code = s.v.voicevoxSynthesizerCreateAudioQueryFromKana(s.synthesizer, text, styleID, uintptr(unsafe.Pointer(&ptr)))
	} else {
		code = s.v.voicevoxSynthesizerCreateAudioQuery(s.synthesizer, text, styleID, uintptr(unsafe.Pointer(&ptr)))
	}
	if code != VOICEVOX_RESULT_OK {
		err = s.v.newError(code)
		return
	}
	defer s.v.voicevoxJsonFree(uintptr(unsafe.Pointer(ptr)))
	j := strings.GoString(ptr)
	err = json.Unmarshal([]byte(j), &a)
	return
}

// 音声合成用のクエリを生成する
func (s *Synthesizer) CreateAudioQuery(text string, styleID StyleId) (a AudioQuery, err error) {
	return s.createAudioQuery(text, styleID, false)
}

// 音声合成用のクエリを生成する(AquesTalk風記法)
func (s *Synthesizer) CreateAudioQueryFromKana(text string, styleID StyleId) (a AudioQuery, err error) {
	return s.createAudioQuery(text, styleID, true)
}
