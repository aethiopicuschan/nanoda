package minimum

import "github.com/ebitengine/purego"

type MinimumCore struct {
	voicevoxGetVersion func() string
}

func NewCore(lib uintptr) (c *MinimumCore, err error) {
	c = &MinimumCore{}
	purego.RegisterLibFunc(&c.voicevoxGetVersion, lib, "voicevox_get_version")
	return
}

func (c *MinimumCore) GetVersion() string {
	return c.voicevoxGetVersion()
}
