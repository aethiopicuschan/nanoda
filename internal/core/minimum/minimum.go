package minimum

import "github.com/ebitengine/purego"

type Core struct {
	core               uintptr
	voicevoxGetVersion func() string
}

func New(core uintptr) (c Core) {
	c.core = core
	purego.RegisterLibFunc(&c.voicevoxGetVersion, c.core, "voicevox_get_version")
	return
}

func (c *Core) GetVersion() string {
	return c.voicevoxGetVersion()
}
