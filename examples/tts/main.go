package main

import (
	"io"
	"os"

	"github.com/aethiopicuschan/nanoda"
)

func main() {
	v, _ := nanoda.NewVoicevox("voicevox_core/libvoicevox_core.dylib", "voicevox_core/open_jtalk_dic_utf_8-1.11", "voicevox_core/model")
	s, _ := v.NewSynthesizer()
	s.LoadModelsFromStyleId(3)
	wav, _ := s.Tts("ずんだもんなのだ！", 3)
	defer wav.Close()
	f, _ := os.Create("output.wav")
	defer f.Close()
	io.Copy(f, wav)
}
