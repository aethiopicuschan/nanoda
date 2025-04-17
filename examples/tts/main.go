package main

import (
	"log"

	"github.com/aethiopicuschan/nanoda/v2"
)

func main() {
	v, _ := nanoda.New("libvoicevox_core.dylib", "open_jtalk_dic_utf_8-1.11", "models")
	log.Println(v.Version())
}
