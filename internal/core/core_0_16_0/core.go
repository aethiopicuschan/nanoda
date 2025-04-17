package core_0_16_0

import "github.com/ebitengine/purego"

type Core struct {
	core               uintptr
	voicevoxGetVersion func() string
}

func New(core uintptr) (c *Core) {
	c = &Core{
		core:               core,
		voicevoxGetVersion: func() string { return "0.16.0" },
	}
	purego.RegisterLibFunc(&c.voicevoxGetVersion, c.core, "voicevox_get_version")
	return
}

func (c *Core) GetVersion() string {
	return c.voicevoxGetVersion()
}
