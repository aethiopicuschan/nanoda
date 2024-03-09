package core_0_15_0

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/aethiopicuschan/nanoda/internal/strings"
	"github.com/aethiopicuschan/nanoda/model"
	"github.com/aethiopicuschan/nanoda/types"
	"github.com/ebitengine/purego"
)

type Core struct {
	// 設定
	accelerationMode int32
	cpuNumThreads    uint16
	// コアの関数群
	// voicevox_audio_query
	// voicevox_audio_query_json_free
	// voicevox_decode
	// voicevox_decode_data_free
	voicevoxErrorResultToMessage func(code types.ResultCode) string
	// voicevox_finalize
	voicevoxGetMetasJson func() string
	// voicevox_get_supported_devices_json
	voicevoxGetVersion    func() string
	voicevoxInitialize    func(uintptr) types.ResultCode
	voicevoxIsGpuMode     func() bool
	voicevoxIsModelLoaded func(id uint32) bool
	voicevoxLoadModel     func(id uint32) types.ResultCode
	// voicevox_make_default_audio_query_options
	// voicevox_make_default_initialize_options
	// voicevox_make_default_synthesis_options
	// voicevox_make_default_tts_options
	// voicevox_predict_duration
	// voicevox_predict_duration_data_free
	// voicevox_predict_intonation
	// voicevox_predict_intonation_data_free
	// voicevox_synthesis
	// voicevox_tts
	// voicevox_wav_free
}

func NewCore(lib uintptr, openJtalkPath string, accelerationMode int, cpuNumThreads int) (c *Core, err error) {
	// オプションを適用する
	c = &Core{
		accelerationMode: int32(accelerationMode),
		cpuNumThreads:    uint16(cpuNumThreads),
	}

	// 関数群の紐付け
	purego.RegisterLibFunc(&c.voicevoxErrorResultToMessage, lib, "voicevox_error_result_to_message")
	purego.RegisterLibFunc(&c.voicevoxGetMetasJson, lib, "voicevox_get_metas_json")
	purego.RegisterLibFunc(&c.voicevoxGetVersion, lib, "voicevox_get_version")
	purego.RegisterLibFunc(&c.voicevoxInitialize, lib, "voicevox_initialize")
	purego.RegisterLibFunc(&c.voicevoxIsGpuMode, lib, "voicevox_is_gpu_mode")
	purego.RegisterLibFunc(&c.voicevoxIsModelLoaded, lib, "voicevox_is_model_loaded")
	purego.RegisterLibFunc(&c.voicevoxLoadModel, lib, "voicevox_load_model")

	// 初期化
	initializeOptions := VoicevoxInitializeOptions{
		accelerationMode: c.accelerationMode,
		cpuNumThreads:    c.cpuNumThreads,
		loadAllModels:    false,
		openJtalkDictDir: strings.CString(openJtalkPath),
	}
	code := c.voicevoxInitialize(*(*uintptr)(unsafe.Pointer(&initializeOptions)))
	if code != types.VOICEVOX_RESULT_OK {
		err = fmt.Errorf(c.voicevoxErrorResultToMessage(code))
		return
	}

	return
}

func (c *Core) ErrorMessageFrom(code types.ResultCode) string {
	return c.voicevoxErrorResultToMessage(code)
}

func (c *Core) GetMetas() (metas []model.Meta, err error) {
	var localMetas []Meta
	metasJson := c.voicevoxGetMetasJson()
	err = json.Unmarshal([]byte(metasJson), &localMetas)
	if err != nil {
		return
	}
	for _, m := range localMetas {
		metas = append(metas, m.ToCommon())
	}
	return
}

func (c *Core) GetVersion() string {
	return c.voicevoxGetVersion()
}

func (c *Core) IsGpuMode() bool {
	return c.voicevoxIsGpuMode()
}

func (c *Core) IsModelLoaded(id int) bool {
	return c.voicevoxIsModelLoaded(uint32(id))
}

func (c *Core) LoadModel(id int) (err error) {
	code := c.voicevoxLoadModel(uint32(id))
	if code != types.VOICEVOX_RESULT_OK {
		err = fmt.Errorf(c.voicevoxErrorResultToMessage(code))
	}
	return
}
