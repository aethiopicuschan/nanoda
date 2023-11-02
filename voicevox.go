package nanoda

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/aethiopicuschan/nanoda/internal/strings"
	"github.com/ebitengine/purego"
)

// 各種関数やポインタなどを保持する構造体
type Voicevox struct {
	core        uintptr
	openJtalkRc uintptr
	vvms        []voicevoxVoiceModel
	// 生(?)の関数群
	voicevoxCreateSupportedDevicesJson             func(uintptr) ResultCode
	voicevoxErrorResultToMessage                   func(ResultCode) string
	voicevoxGetVersion                             func() string
	voicevoxJsonFree                               func(uintptr)
	voicevoxSynthesizerCreateAccentPhrases         func(uintptr, string, StyleId, uintptr) ResultCode
	voicevoxSynthesizerCreateAccentPhrasesFromKana func(uintptr, string, StyleId, uintptr) ResultCode
	voicevoxSynthesizerCreateAudioQuery            func(uintptr, string, StyleId, uintptr) ResultCode
	voicevoxSynthesizerCreateAudioQueryFromKana    func(uintptr, string, StyleId, uintptr) ResultCode
	voicevoxSynthesizerCreateMetasJson             func(uintptr) uintptr
	voicevoxOpenJtalkRcDelete                      func(uintptr)
	voicevoxOpenJtalkRcNew                         func(string, uintptr) ResultCode
	voicevoxOpenJtalkRcUseUserDict                 func(uintptr, uintptr) ResultCode
	voicevoxSynthesizerDelete                      func(uintptr)
	voicevoxSynthesizerIsGpuMode                   func(uintptr) bool
	voicevoxSynthesizerIsLoadedVoiceModel          func(uintptr, string) bool
	voicevoxSynthesizerLoadVoiceModel              func(uintptr, uintptr) ResultCode
	voicevoxSynthesizerNewWithInitialize           func(uintptr, uintptr, uintptr) ResultCode
	voicevoxSynthesizerReplaceMoraData             func(uintptr, string, StyleId, uintptr) ResultCode
	voicevoxSynthesizerReplaceMoraPitch            func(uintptr, string, StyleId, uintptr) ResultCode
	voicevoxSynthesizerReplacePhonemeLength        func(uintptr, string, StyleId, uintptr) ResultCode
	voicevoxSynthesizerSynthesis                   func(uintptr, string, StyleId, uintptr, uintptr, uintptr) ResultCode
	voicevoxSynthesizerTts                         func(uintptr, string, StyleId, uintptr, uintptr, uintptr) ResultCode
	voicevoxSynthesizerTtsFromKana                 func(uintptr, string, StyleId, uintptr, uintptr, uintptr) ResultCode
	voicevoxSynthesizerUnloadVoiceModel            func(uintptr, string) ResultCode
	voicevoxUserDictAddWord                        func(uintptr, uintptr, uintptr) ResultCode
	voicevoxUserDictDelete                         func(uintptr)
	voicevoxUserDictImport                         func(uintptr, uintptr) ResultCode
	voicevoxUserDictLoad                           func(uintptr, string) ResultCode
	voicevoxUserDictNew                            func() uintptr
	voicevoxUserDictRemoveWord                     func(uintptr, uintptr) ResultCode
	voicevoxUserDictSave                           func(uintptr, string) ResultCode
	voicevoxUserDictToJson                         func(uintptr, uintptr) ResultCode
	voicevoxUserDictUpdateWord                     func(uintptr, uintptr, uintptr) ResultCode
	voicevoxVoiceModelDelete                       func(uintptr)
	voicevoxVoiceModelGetMetasJson                 func(uintptr) string
	voicevoxVoiceModelId                           func(uintptr) string
	voicevoxVoiceModelNewFromPath                  func(string, uintptr) ResultCode
	voicevoxWavFree                                func(uintptr)
}

// 必要なパスを引数に取り、Voicevoxのインスタンスを生成する
func NewVoicevox(corePath string, openJtalkPath string, modelPath string) (v *Voicevox, err error) {
	c, err := purego.Dlopen(corePath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return
	}
	v = &Voicevox{
		core: c,
	}
	// 関数の紐付けを行う
	v.register()
	// OpenJtalkRcを構築する
	code := v.voicevoxOpenJtalkRcNew(openJtalkPath, uintptr(unsafe.Pointer(&v.openJtalkRc)))
	if code != VOICEVOX_RESULT_OK {
		err = v.newError(code)
		return
	}
	// VVMファイルを読み込む
	v.vvms, err = v.loadVVMs(modelPath)
	// エラーがあった場合、もろもろを破棄する
	if err != nil {
		v.voicevoxOpenJtalkRcDelete(v.openJtalkRc)
		v.deleteVVMs()
	}
	return
}

// 利用可能なデバイスの情報を取得する
func (v *Voicevox) SupportedDevices() (sd SupportedDevices, err error) {
	var ptr *byte
	code := v.voicevoxCreateSupportedDevicesJson(uintptr(unsafe.Pointer(&ptr)))
	defer v.voicevoxJsonFree(uintptr(unsafe.Pointer(ptr)))
	if code != VOICEVOX_RESULT_OK {
		err = v.newError(code)
		return
	}
	j := strings.GoString(ptr)
	err = json.Unmarshal([]byte(j), &sd)
	return
}

// 結果コードからメッセージを取得する
func (v *Voicevox) GetMessageFromResult(code ResultCode) string {
	return v.voicevoxErrorResultToMessage(code)
}

// voicevoxのバージョンを取得する
func (v *Voicevox) GetVersion() string {
	return v.voicevoxGetVersion()
}

// メタ情報の一覧を取得する
func (v *Voicevox) GetMetas() []Meta {
	return v.sortMetas(v.getMetas())
}

/*
スタイルの一覧を取得する
ここで取得されるスタイルのNameは「ずんだもん(あまあま)」のようなフォーマットになる
*/
func (v *Voicevox) GetStyles() []Style {
	metas := v.sortMetas(v.getMetas())
	styles := []Style{}
	for _, meta := range metas {
		for _, style := range meta.Styles {
			style.Name = fmt.Sprintf("%s(%s)", meta.Name, style.Name)
			styles = append(styles, style)
		}
	}
	return styles
}
