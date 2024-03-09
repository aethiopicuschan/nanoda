package core_0_15_0

import (
	"github.com/aethiopicuschan/nanoda/constant"
	"github.com/ebitengine/purego"
)

type Core struct {
	// 設定
	accelerationMode constant.AccelerationMode
	cpuNumThreads    int
	// コアの関数群
	voicevoxGetVersion           func() string
	voicevoxErrorResultToMessage func(code constant.ResultCode) string
}

func NewCore(lib uintptr, accelerationMode constant.AccelerationMode, cpuNumThreads int) (c *Core, err error) {
	c = &Core{
		accelerationMode: accelerationMode,
		cpuNumThreads:    cpuNumThreads,
	}
	// 関数群の紐付け
	purego.RegisterLibFunc(&c.voicevoxGetVersion, lib, "voicevox_get_version")
	purego.RegisterLibFunc(&c.voicevoxErrorResultToMessage, lib, "voicevox_error_result_to_message")
	return
}

func (c *Core) ErrorMessageFrom(code constant.ResultCode) string {
	return c.voicevoxErrorResultToMessage(code)
}

func (c *Core) GetVersion() string {
	return c.voicevoxGetVersion()
}
