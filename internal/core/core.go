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
	AccelerationMode int
	CpuNumThreads    int
}

func NewCore(lib string, openJtalkPath string, o *Option) (c Core, err error) {
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

	// バージョンによってCoreを選択して返す
	switch version {
	case "0.15.0":
		c, err = core_0_15_0.NewCore(l, openJtalkPath, o.AccelerationMode, o.CpuNumThreads)
	default:
		err = fmt.Errorf("unsupported version: %s", version)
	}
	return
}
