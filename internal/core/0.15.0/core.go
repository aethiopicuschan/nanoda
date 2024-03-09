package core_0_15_0

import (
	"github.com/ebitengine/purego"
)

type Core struct {
	voicevoxGetVersion           func() string
	voicevoxErrorResultToMessage func(code int) string
}

func NewCore(lib uintptr) (c *Core, err error) {
	c = &Core{}
	purego.RegisterLibFunc(&c.voicevoxGetVersion, lib, "voicevox_get_version")
	purego.RegisterLibFunc(&c.voicevoxErrorResultToMessage, lib, "voicevox_error_result_to_message")
	return
}

func (c *Core) ErrorMessageFrom(code int) string {
	return c.voicevoxErrorResultToMessage(code)
}

func (c *Core) GetVersion() string {
	return c.voicevoxGetVersion()
}
