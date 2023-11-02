package nanoda

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/aethiopicuschan/nanoda/internal/strings"
)

// アクセント句の再生成時に指定するオプション
type replaceOptions struct {
	replaceMoraPitch     bool
	replacePhonemeLength bool
}

// 音高を生成しなおす
func withReplaceMoraPitch() func(*replaceOptions) {
	return func(o *replaceOptions) {
		o.replaceMoraPitch = true
	}
}

// 音素長を生成しなおす
func withReplacePhonemeLength() func(*replaceOptions) {
	return func(o *replaceOptions) {
		o.replacePhonemeLength = true
	}
}

// アクセント句の配列を指定されたスタイルで再生成する
func (s *Synthesizer) replace(ap []AccentPhrase, styleID StyleId, options ...func(*replaceOptions)) (a []AccentPhrase, err error) {
	opt := replaceOptions{}
	for _, o := range options {
		o(&opt)
	}
	var ptr *byte
	var code ResultCode
	jB, err := json.Marshal(ap)
	if err != nil {
		return
	}
	j := string(jB)
	switch {
	case opt.replaceMoraPitch && opt.replacePhonemeLength:
		code = s.v.voicevoxSynthesizerReplaceMoraData(s.synthesizer, j, styleID, uintptr(unsafe.Pointer(&ptr)))
	case opt.replaceMoraPitch:
		code = s.v.voicevoxSynthesizerReplaceMoraPitch(s.synthesizer, j, styleID, uintptr(unsafe.Pointer(&ptr)))
	case opt.replacePhonemeLength:
		code = s.v.voicevoxSynthesizerReplacePhonemeLength(s.synthesizer, j, styleID, uintptr(unsafe.Pointer(&ptr)))
	default:
		err = fmt.Errorf("options are required")
		return
	}
	if code != VOICEVOX_RESULT_OK {
		err = s.v.newError(code)
		return
	}
	defer s.v.voicevoxJsonFree(uintptr(unsafe.Pointer(ptr)))
	j = strings.GoString(ptr)
	err = json.Unmarshal([]byte(j), &a)
	return
}

// アクセント句の配列を指定されたスタイルで再生成する
func (s *Synthesizer) Replace(ap []AccentPhrase, styleID StyleId) (a []AccentPhrase, err error) {
	return s.replace(ap, styleID, withReplaceMoraPitch(), withReplacePhonemeLength())
}

// アクセント句の配列を指定されたスタイルで再生成する(音高のみ)
func (s *Synthesizer) ReplaceOnlyMoraPitch(ap []AccentPhrase, styleID StyleId) (a []AccentPhrase, err error) {
	return s.replace(ap, styleID, withReplaceMoraPitch())
}

// アクセント句の配列を指定されたスタイルで再生成する(音素長のみ)
func (s *Synthesizer) ReplaceOnlyPhonemeLength(ap []AccentPhrase, styleID StyleId) (a []AccentPhrase, err error) {
	return s.replace(ap, styleID, withReplacePhonemeLength())
}
