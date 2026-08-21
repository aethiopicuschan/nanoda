// Command tts is a minimal end-to-end smoke test for nanoda's 0.16.x
// support: load the library, build a synthesizer, load one voice model,
// and write out a WAV file.
//
// It expects a voicevox_core 0.16.x runtime laid out as:
//
//	voicevox_core.dll (or libvoicevox_core.so/.dylib) — this directory
//	voicevox_onnxruntime.dll (or the .so/.dylib equivalent) — this
//	  directory too, alongside voicevox_core's own library file, *not*
//	  under runtime/. voicevox_core loads it by bare filename via the
//	  OS's default shared-library search path, which (on Windows, with
//	  the default Safe DLL Search Mode) checks the running process's own
//	  directory and the current directory, but never subdirectories of
//	  either — one level too deep to find it under runtime/. Everything
//	  else below is opened via an explicit path from Go instead, so
//	  nesting it doesn't matter.
//	runtime/dict/open_jtalk_dic_utf_8-1.11/
//	runtime/models/*.vvm
//
// See nanoda's README "必要なファイルについて" for how to obtain a
// voicevox_core runtime in the first place.
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/aethiopicuschan/nanoda/v2"
)

func main() {
	corePath := corePathForOS()
	v, err := nanoda.New(corePath, "runtime/dict/open_jtalk_dic_utf_8-1.11", "runtime/models")
	if err != nil {
		log.Fatalf("nanoda.New: %v", err)
	}
	log.Println("voicevox_core version:", v.Version())

	s, err := v.NewSynthesizer()
	if err != nil {
		log.Fatalf("NewSynthesizer: %v", err)
	}
	defer s.Close()

	if err := s.LoadAllModels(); err != nil {
		log.Fatalf("LoadAllModels: %v", err)
	}

	// Style 3 is ずんだもん(ノーマル) in the official model numbering
	// (see model 0.vvm's own metas).
	wav, err := s.Tts("ボイスボックス コア ゼロテンイチロク、動作確認なのだ", 3)
	if err != nil {
		log.Fatalf("Tts: %v", err)
	}
	defer wav.Close()

	f, err := os.Create("output.wav")
	if err != nil {
		log.Fatalf("os.Create: %v", err)
	}
	defer f.Close()
	if _, err := f.ReadFrom(wav); err != nil {
		log.Fatalf("write output.wav: %v", err)
	}

	log.Println("wrote output.wav")
}

func corePathForOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "libvoicevox_core.dylib"
	case "linux":
		return "libvoicevox_core.so"
	default:
		return "voicevox_core.dll"
	}
}
