package core_0_17_0

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"unsafe"

	"github.com/aethiopicuschan/nanoda/v2/internal/strings"
)

const (
	defaultSamplingRate = 24000
	// defaultSamplingRate / VOICEVOX_FRAME_RATE (93.75).
	waveSamplesPerFrame = 256
	// voicevox_core's own streaming_synthesis default, about 0.3s per chunk.
	segmentFrames = 28
)

// audioFeatureSlice mirrors VoicevoxAudioFeatureSlice. Unlike the options
// structs it is passed to the registered func as a real struct; see the
// package doc comment.
type audioFeatureSlice struct {
	audioFeature         uintptr
	startOffset          uintptr
	end                  uintptr
	endIsInclusive       bool
	clamp                bool
	allowStartExceedsEnd bool
}

// synthesisOptions mirrors VoicevoxSynthesisOptions (a single bool, 1
// byte).
type synthesisOptions struct {
	enableInterrogativeUpspeak bool
}

// SupportsStreaming reports whether the loaded library has the 0.17.1+
// functions TtsStream needs.
func (s *Synthesizer) SupportsStreaming() bool {
	return s.c.supportsStreaming
}

// TtsStream synthesizes text like Tts, but renders the audio lazily as the
// returned WAV stream is read: the 44-byte header comes first, then each
// Read that runs out of buffered PCM renders the next segmentFrames
// frames. The stream must be closed before the Synthesizer.
func (s *Synthesizer) TtsStream(text string, styleID uint32) (io.ReadCloser, error) {
	if !s.c.supportsStreaming {
		return nil, errors.New("nanoda: streaming synthesis requires voicevox_core 0.17.1 or later")
	}

	var queryPtr *byte
	if code := s.c.voicevoxSynthesizerCreateAudioQuery(s.synthesizer, text, styleID, &queryPtr); code != 0 {
		return nil, s.c.newError("synthesizer_create_audio_query", code)
	}
	defer s.c.voicevoxStringFree(uintptr(unsafe.Pointer(queryPtr)))
	// Aliases the library's buffer, so it must not be used after the
	// deferred free above.
	query := strings.GoString(queryPtr)

	var format struct {
		OutputSamplingRate uint32 `json:"outputSamplingRate"`
		OutputStereo       bool   `json:"outputStereo"`
	}
	if err := json.Unmarshal([]byte(query), &format); err != nil {
		return nil, fmt.Errorf("nanoda: parse audio query: %w", err)
	}

	// Same as Tts, so both produce identical audio for the same input.
	opts := synthesisOptions{enableInterrogativeUpspeak: false}
	var feature uintptr
	if code := s.c.voicevoxSynthesizerCreateAudioFeature(s.synthesizer, query, styleID, *(*uintptr)(unsafe.Pointer(&opts)), &feature); code != 0 {
		return nil, s.c.newError("synthesizer_create_audio_feature", code)
	}
	frames := s.c.voicevoxAudioFeatureLen(feature)

	header := wavHeader(frames, format.OutputSamplingRate, format.OutputStereo)
	return &stream{s: s, feature: feature, frames: frames, pending: header}, nil
}

// wavHeader builds the same 44-byte header voicevox_core's
// streaming_synthesis emits. render's PCM already has volume_scale,
// output_sampling_rate (by repeating each sample) and output_stereo
// applied, so the data size follows directly from the frame count.
func wavHeader(frames uintptr, samplingRate uint32, stereo bool) []byte {
	channels := uint16(1)
	if stereo {
		channels = 2
	}
	blockAlign := channels * 2
	repeat := samplingRate / defaultSamplingRate * uint32(channels)
	pcmLen := uint32(frames) * waveSamplesPerFrame * repeat * 2

	h := make([]byte, 0, 44)
	h = append(h, "RIFF"...)
	h = binary.LittleEndian.AppendUint32(h, 36+pcmLen)
	h = append(h, "WAVEfmt "...)
	h = binary.LittleEndian.AppendUint32(h, 16)
	h = binary.LittleEndian.AppendUint16(h, 1) // linear PCM
	h = binary.LittleEndian.AppendUint16(h, channels)
	h = binary.LittleEndian.AppendUint32(h, samplingRate)
	h = binary.LittleEndian.AppendUint32(h, samplingRate*uint32(blockAlign))
	h = binary.LittleEndian.AppendUint16(h, blockAlign)
	h = binary.LittleEndian.AppendUint16(h, 16)
	h = append(h, "data"...)
	h = binary.LittleEndian.AppendUint32(h, pcmLen)
	return h
}

// stream is the io.ReadCloser TtsStream returns. It owns the native
// VoicevoxAudioFeature until Close.
type stream struct {
	s       *Synthesizer
	feature uintptr
	frames  uintptr
	cursor  uintptr
	// pending is the unread part of buf (or the header, initially).
	pending []byte
	buf     []byte
}

func (st *stream) Read(p []byte) (int, error) {
	if st.feature == 0 {
		return 0, errors.New("nanoda: read from closed stream")
	}
	for len(st.pending) == 0 {
		if st.cursor >= st.frames {
			return 0, io.EOF
		}
		if err := st.renderNext(); err != nil {
			return 0, err
		}
	}
	n := copy(p, st.pending)
	st.pending = st.pending[n:]
	return n, nil
}

// renderNext renders [cursor, cursor+segmentFrames) the same way
// voicevox_core's own SynthesisStreamInner slices the feature: a plain
// half-open range, no clamping. Margins around the segment are handled
// inside voicevox_synthesizer_render.
func (st *stream) renderNext() error {
	c := st.s.c
	end := min(st.cursor+segmentFrames, st.frames)
	var slice audioFeatureSlice
	if code := c.voicevoxAudioFeatureSlice(st.feature, st.cursor, end, false, false, false, &slice); code != 0 {
		return c.newError("audio_feature_slice", code)
	}
	var pcmLen uintptr
	var pcmPtr *byte
	if code := c.voicevoxSynthesizerRender(st.s.synthesizer, slice, &pcmLen, &pcmPtr); code != 0 {
		return c.newError("synthesizer_render", code)
	}
	st.cursor = end
	// A zero-length result points at voicevox_empty_bytes, which
	// voicevox_bytes_free warns about instead of freeing.
	if pcmLen == 0 {
		return nil
	}
	// Copied so the native buffer is freed now rather than tracked until
	// the caller has read it all.
	st.buf = append(st.buf[:0], unsafe.Slice(pcmPtr, pcmLen)...)
	c.bytesFree(uintptr(unsafe.Pointer(pcmPtr)))
	st.pending = st.buf
	return nil
}

func (st *stream) Close() error {
	if st.feature != 0 {
		st.s.c.voicevoxAudioFeatureDelete(st.feature)
		st.feature = 0
	}
	return nil
}
