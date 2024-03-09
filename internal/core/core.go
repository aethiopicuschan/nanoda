package core

import (
	"fmt"

	core_0_15_0 "github.com/aethiopicuschan/nanoda/internal/core/0.15.0"
	"github.com/aethiopicuschan/nanoda/internal/core/minimum"
)

type Core interface {
	GetVersion() string
	ErrorMessageFrom(code int) string
}

func NewCore(lib string) (c Core, err error) {
	l, err := openLibrary(lib)
	if err != nil {
		return
	}
	mc, err := minimum.NewCore(l)
	if err != nil {
		return
	}
	version := mc.GetVersion()
	switch version {
	case "0.15.0":
		c, err = core_0_15_0.NewCore(l)
	default:
		err = fmt.Errorf("unsupported version: %s", version)
	}
	return
}
