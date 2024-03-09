package nanoda

import "github.com/aethiopicuschan/nanoda/internal/core"

type Voicevox struct {
	core core.Core
}

func NewVoicevox(corePath string, openJtalkPath string) (v *Voicevox, err error) {
	v = &Voicevox{}
	v.core, err = core.NewCore(corePath)
	if err != nil {
		return
	}
	return
}

// voicevox coreのバージョンを取得する
func (v *Voicevox) GetVersion() string {
	return v.core.GetVersion()
}
