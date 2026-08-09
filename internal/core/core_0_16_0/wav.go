package core_0_16_0

import (
	"bytes"
	"io"
)

// wav wraps voicevox_synthesizer_tts's output buffer as an io.ReadCloser,
// freeing it (via voicevox_wav_free) on Close rather than earlier — main's
// 0.15 implementation used the same shape (its own wav.go).
type wav struct {
	reader io.Reader
	close  func() error
}

func newWav(raw []byte, close func() error) *wav {
	return &wav{reader: bytes.NewReader(raw), close: close}
}

func (w *wav) Read(p []byte) (n int, err error) {
	return w.reader.Read(p)
}

func (w *wav) Close() error {
	return w.close()
}
