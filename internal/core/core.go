package core

import (
	"github.com/aethiopicuschan/nanoda/v2/internal/core/core_0_16_0"
	"github.com/aethiopicuschan/nanoda/v2/internal/core/minimum"
)

type Core interface {
	GetVersion() string
}

func New(core uintptr) (c Core) {
	mc := minimum.New(core)
	switch {
	case mc.GetVersion() == "0.16.0":
		return core_0_16_0.New(core)
	case mc.GetVersion() == "":
		return nil
	}
	return
}
