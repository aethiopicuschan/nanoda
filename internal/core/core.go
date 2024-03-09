package core

import (
	"fmt"

	"github.com/aethiopicuschan/nanoda/constant"
	core_0_15_0 "github.com/aethiopicuschan/nanoda/internal/core/0.15.0"
	"github.com/aethiopicuschan/nanoda/internal/core/minimum"
	"github.com/aethiopicuschan/nanoda/model"
)

// バージョンに依存しない抽象化されたコアのインターフェース
type Core interface {
	GetVersion() string
	IsGpuMode() bool
	ErrorMessageFrom(code constant.ResultCode) string
	GetMetas() (metas []model.Meta, err error)
}

// オプション
type Option struct {
	AccelerationMode int
	CpuNumThreads    int
}

// 対象のlibからバージョンを取得して対応するコアを返す
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

	// バージョンによってコアを選択して返す
	switch version {
	case "0.15.0":
		c, err = core_0_15_0.NewCore(l, openJtalkPath, o.AccelerationMode, o.CpuNumThreads)
	default:
		err = fmt.Errorf("unsupported version: %s", version)
	}
	return
}
