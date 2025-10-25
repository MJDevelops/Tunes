package audio

import (
	"errors"
	"slices"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/mjdevelops/tunes/internal/pkg/os"
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

func NewDecoder(file string) (Decoder, error) {
	ext := os.GetFileExtension(file)
	if err := IsSupportedFormat(ext); err != nil {
		return nil, err
	}

	switch ext {
	case ".wav":
		return NewWavDecoder(file)
	default:
		return NewTagDecoder(file)
	}
}

func (ad *AudioFile) Duration() time.Duration {
	return ad.buffer.Format().SampleRate.D(ad.buffer.Len())
}

// TODO: Implement this
func (ad *AudioFile) Play()   {}
func (ad *AudioFile) Pause()  {}
func (ad *AudioFile) Resume() {}

// IsSupportedFormat reports whether the provided format is supported
// in the scope of audio decoding. If the format is not supported this
// will return an error of type ErrUnsupported.
func IsSupportedFormat(format string) error {
	if !slices.Contains(supportedFormats, format) {
		return ErrUnsupported
	}
	return nil
}
