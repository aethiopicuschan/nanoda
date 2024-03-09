package core

import (
	"fmt"

	"github.com/aethiopicuschan/nanoda/constant"
	core_0_15_0 "github.com/aethiopicuschan/nanoda/internal/core/0.15.0"
	"github.com/aethiopicuschan/nanoda/internal/core/minimum"
)

type Core interface {
	GetVersion() string
	ErrorMessageFrom(code constant.ResultCode) string
}

type Option struct {
	accelerationMode constant.AccelerationMode
	cpuNumThreads    int
}

// ハードウェアアクセラレーションモードを設定する
func WithAccelerationMode(mode constant.AccelerationMode) func(*Option) {
	return func(o *Option) {
		o.accelerationMode = mode
	}
}

// CPU利用数を設定する 0の場合は環境に合わせてCPUが利用される
func WithCpuNumThreads(num int) func(*Option) {
	return func(o *Option) {
		o.cpuNumThreads = num
	}
}

func NewCore(lib string, options ...func(*Option)) (c Core, err error) {
	// ライブラリを開く
	l, err := openLibrary(lib)
	if err != nil {
		return
	}

	// バージョンを取得する
	mc, err := minimum.NewCore(l)
	if err != nil {
		return
	}
	version := mc.GetVersion()

	// デフォルトオプション
	o := &Option{
		accelerationMode: constant.ACCELERATION_MODE_AUTO,
		cpuNumThreads:    0,
	}
	// オプションを適用する
	for _, option := range options {
		option(o)
	}

	// バージョンによってCoreを選択して返す
	switch version {
	case "0.15.0":
		c, err = core_0_15_0.NewCore(l, o.accelerationMode, o.cpuNumThreads)
	default:
		err = fmt.Errorf("unsupported version: %s", version)
	}
	return
}
