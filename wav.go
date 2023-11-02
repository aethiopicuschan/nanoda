package nanoda

import (
	"bytes"
	"io"
)

// 出力されたwavファイルを表す構造体
// io.ReadCloserを実装している
type Wav struct {
	reader io.Reader
	close  func() error
}

func newWav(raw []byte, close func() error) *Wav {
	return &Wav{
		reader: bytes.NewReader(raw),
		close:  close,
	}
}

func (w *Wav) Read(p []byte) (n int, err error) {
	return w.reader.Read(p)
}

func (w *Wav) Close() error {
	return w.close()
}
