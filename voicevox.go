package nanoda

import (
	"github.com/aethiopicuschan/nanoda/v2/internal/core"
	"github.com/aethiopicuschan/nanoda/v2/internal/open"
)

type Voicevox struct {
	core core.Core
}

func New(corePath string, openJtalkPath string, modelPath string) (v *Voicevox, err error) {
	c, err := open.Open(corePath)
	if err != nil {
		return
	}
	core := core.New(c)
	if core == nil {
		err = ErrUnsupportedVersion
		return
	}
	v = &Voicevox{
		core: core,
	}
	return
}

func (v *Voicevox) Version() string {
	return v.core.GetVersion()
}
