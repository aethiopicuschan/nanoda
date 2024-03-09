package main

import (
	"log"

	"github.com/aethiopicuschan/nanoda"
)

func main() {
	v, err := nanoda.NewVoicevox("voicevox_core-0.15.0/libvoicevox_core.dylib", "voicevox_core-0.15.0/open_jtalk_dic_utf_8-1.11")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v.IsGpuMode())
}
