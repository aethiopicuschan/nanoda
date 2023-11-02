package nanoda

import (
	"encoding/json"
	"io"
	"unsafe"

	"github.com/aethiopicuschan/nanoda/internal/strings"
)

// ハードウェアアクセラレーションモードを設定する設定値
type AccelerationMode int32

const (
	ACCELERATION_MODE_AUTO AccelerationMode = iota // 実行環境に合った適切なハードウェアアクセラレーションモードを選択する
	ACCELERATION_MODE_CPU                          // ハードウェアアクセラレーションモードを"CPU"に設定する
	ACCELERATION_MODE_GPU                          // ハードウェアアクセラレーションモードを"GPU"に設定する
)

// シンセナイザの作成時に指定するオプション
type SynthesizerOption struct {
	accelerationMode AccelerationMode
	cpuNumThreads    uint16
}

// ハードウェアアクセラレーションモードを設定する
func WithAccelerationMode(mode AccelerationMode) func(*SynthesizerOption) {
	return func(o *SynthesizerOption) {
		o.accelerationMode = mode
	}
}

// CPU利用数を設定する 0の場合は環境に合わせてCPUが利用される
func WithCpuNumThreads(num uint16) func(*SynthesizerOption) {
	return func(o *SynthesizerOption) {
		o.cpuNumThreads = num
	}
}

// 音声シンセナイザ
// Voicevox.NewSynthesizerで作成する
type Synthesizer struct {
	v           *Voicevox
	synthesizer uintptr
}

// シンセナイザを作成する
func (v *Voicevox) NewSynthesizer(options ...func(*SynthesizerOption)) (s Synthesizer, err error) {
	s.v = v
	opt := SynthesizerOption{
		accelerationMode: ACCELERATION_MODE_AUTO,
		cpuNumThreads:    0,
	}
	for _, o := range options {
		o(&opt)
	}
	code := v.voicevoxSynthesizerNewWithInitialize(v.openJtalkRc, *(*uintptr)(unsafe.Pointer(&opt)), uintptr(unsafe.Pointer(&s.synthesizer)))
	if code != VOICEVOX_RESULT_OK {
		err = v.newError(code)
	}
	return
}

// このSynthesizerを閉じる
func (s *Synthesizer) Close() {
	s.v.voicevoxSynthesizerDelete(s.synthesizer)
}

// 現在読み込んでいる音声モデルのメタ情報を取得する
func (s *Synthesizer) GetMetas() (metas []Meta, err error) {
	ptr := s.v.voicevoxSynthesizerCreateMetasJson(s.synthesizer)
	defer s.v.voicevoxJsonFree(ptr)
	j := strings.GoStringFromUintptr(ptr)
	ms := []meta{}
	if err = json.Unmarshal([]byte(j), &ms); err != nil {
		return
	}
	metas = s.v.sortMetas(ms)
	return
}

// GPUモードかどうかを判定する
func (s *Synthesizer) IsGpuMode() bool {
	return s.v.voicevoxSynthesizerIsGpuMode(s.synthesizer)
}

// Ttsを用いた音声合成時に指定するオプション
type TtsOptions struct {
	// 疑問文の調整を有効にするかどうか
	enableInterrogativeUpspeak bool
	// AquesTalk風記法を有効にするかどうか
	enableKana bool
}

// 疑問文の調整を有効にする
func WithEnableInterrogativeUpspeak() func(*TtsOptions) {
	return func(o *TtsOptions) {
		o.enableInterrogativeUpspeak = true
	}
}

// AquesTalk風記法を有効にする
func WithEnableKana() func(*TtsOptions) {
	return func(o *TtsOptions) {
		o.enableKana = true
	}
}

// 音声合成を行う
func (s *Synthesizer) Tts(text string, styleID StyleId, options ...func(*TtsOptions)) (io.ReadCloser, error) {
	opt := TtsOptions{
		enableInterrogativeUpspeak: false,
	}
	for _, o := range options {
		o(&opt)
	}
	opt2 := struct {
		enableInterrogativeUpspeak bool
	}{
		opt.enableInterrogativeUpspeak,
	}

	var outputBinarySize uint
	var outputWav *uint8

	var code ResultCode
	if opt.enableKana {
		code = s.v.voicevoxSynthesizerTtsFromKana(s.synthesizer, text, styleID, *(*uintptr)(unsafe.Pointer(&opt2)), uintptr(unsafe.Pointer(&outputBinarySize)), uintptr(unsafe.Pointer(&outputWav)))
	} else {
		code = s.v.voicevoxSynthesizerTts(s.synthesizer, text, styleID, *(*uintptr)(unsafe.Pointer(&opt2)), uintptr(unsafe.Pointer(&outputBinarySize)), uintptr(unsafe.Pointer(&outputWav)))
	}
	if code != VOICEVOX_RESULT_OK {
		return nil, s.v.newError(code)
	}
	raw := unsafe.Slice(outputWav, outputBinarySize)
	wav := newWav(raw, func() error {
		s.v.voicevoxWavFree(uintptr(unsafe.Pointer(outputWav)))
		return nil
	})
	return wav, nil
}

// AudioQueryから音声合成を行う
func (s *Synthesizer) synthesis(aq AudioQuery, styleId StyleId, enableInterrogativeUpspeak bool) (io.ReadCloser, error) {
	opt := struct {
		enableInterrogativeUpspeak bool
	}{
		enableInterrogativeUpspeak,
	}

	jB, err := json.Marshal(aq)
	if err != nil {
		return nil, err
	}
	j := string(jB)

	var outputBinarySize uint
	var outputWav *uint8

	code := s.v.voicevoxSynthesizerSynthesis(s.synthesizer, j, styleId, *(*uintptr)(unsafe.Pointer(&opt)), uintptr(unsafe.Pointer(&outputBinarySize)), uintptr(unsafe.Pointer(&outputWav)))

	if code != VOICEVOX_RESULT_OK {
		return nil, s.v.newError(code)
	}
	raw := unsafe.Slice(outputWav, outputBinarySize)
	wav := newWav(raw, func() error {
		s.v.voicevoxWavFree(uintptr(unsafe.Pointer(outputWav)))
		return nil
	})
	return wav, nil
}

// AudioQueryから音声合成を行う
func (s *Synthesizer) Synthesis(aq AudioQuery, styleId StyleId) (io.ReadCloser, error) {
	return s.synthesis(aq, styleId, true)
}

// AudioQueryから音声合成を行う(疑問文の調整なし)
func (s *Synthesizer) SynthesisWithoutInterrogativeUpspeak(aq AudioQuery, styleId StyleId) (io.ReadCloser, error) {
	return s.synthesis(aq, styleId, false)
}
