package audio

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
)

var supportedFormats = []string{".flac", ".ogg", ".mp3", ".wav"}

type Decoder interface {
	DecodeAudio() (beep.StreamSeekCloser, beep.Format, error)
	ParseMeta() (TrackMeta, error)
}

type AudioFile struct {
	Metadata TrackMeta
	buffer   *beep.Buffer
	format   beep.Format
}

type Queue struct {
	Queue []AudioFile
	mu    sync.Mutex
}

var (
	ErrUnsupported = errors.New("unsupported file format")
)

// NewAudioFile constructs a new AudioFile struct with the provided decoder.
//
// This function returns err != nil if the audio decoding fails.
func NewAudioFile(decoder Decoder) (ad AudioFile, err error) {
	af := AudioFile{}

	var (
		format   beep.Format
		buffer   *beep.Buffer
		streamer beep.StreamSeekCloser
	)

	streamer, format, err = decoder.DecodeAudio()

	if err != nil {
		return af, err
	}

	buffer = beep.NewBuffer(format)
	buffer.Append(streamer)
	af.Metadata, _ = decoder.ParseMeta()
	streamer.Close()

	af.buffer = buffer
	af.format = format

	return af, nil
}

func (ad *AudioFile) Duration() time.Duration {
	return ad.buffer.Format().SampleRate.D(ad.buffer.Len())
}

// IsSupportedFormat reports whether the provided format is supported
// in the scope of audio decoding. If the format is not supported this
// will return an error of type ErrUnsupported.
func IsSupportedFormat(format string) error {
	if !slices.Contains(supportedFormats, format) {
		return ErrUnsupported
	}
	return nil
}
