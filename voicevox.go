package nanoda

import (
	"github.com/aethiopicuschan/nanoda/internal/core"
	"github.com/aethiopicuschan/nanoda/model"
)

type Voicevox struct {
	core core.Core
}

// ハードウェアアクセラレーションモードを設定する
func WithAccelerationMode(mode int) func(*core.Option) {
	return func(o *core.Option) {
		o.AccelerationMode = mode
	}
}

// CPU利用数を設定する 0の場合は環境に合わせてCPUが利用される
func WithCpuNumThreads(num int) func(*core.Option) {
	return func(o *core.Option) {
		o.CpuNumThreads = num
	}
}

func NewVoicevox(corePath string, openJtalkPath string, options ...func(*core.Option)) (v *Voicevox, err error) {
	v = &Voicevox{}

	// デフォルトオプション
	o := &core.Option{
		AccelerationMode: ACCELERATION_MODE_AUTO,
		CpuNumThreads:    0,
	}
	// オプションを適用する
	for _, option := range options {
		option(o)
	}

	v.core, err = core.NewCore(corePath, openJtalkPath, o)
	if err != nil {
		return
	}
	return
}

// メタ情報を取得する
func (v *Voicevox) GetMetas() (metas []model.Meta, err error) {
	return v.core.GetMetas()
}

// voicevox coreのバージョンを取得する
func (v *Voicevox) GetVersion() string {
	return v.core.GetVersion()
}
