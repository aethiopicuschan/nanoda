package nanoda

import "errors"

var ErrUnsupportedVersion = errors.New("unsupported version")

// ErrStreamingUnsupported is returned by (*Synthesizer).TtsStream when the
// loaded voicevox_core is older than 0.17.1.
var ErrStreamingUnsupported = errors.New("streaming synthesis unsupported")
