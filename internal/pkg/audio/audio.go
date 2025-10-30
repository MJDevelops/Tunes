package audio

import (
	"errors"
	"time"

	"github.com/gopxl/beep/v2"
)

type Decoder interface {
	New(path string) (Decoder, error)
	DecodeAudio() (*AudioFile, error)
	ParseMeta() (TrackMeta, error)
}

type AudioFile struct {
	buffer *beep.Buffer
	format beep.Format
}

var supportedFormats map[string]Decoder

var (
	ErrUnsupported = errors.New("unsupported file format")
)

func RegisterDecoder(decoder Decoder, formats ...string) {
	for _, format := range formats {
		supportedFormats[format] = decoder
	}
}

func GetDecoder(format string) (Decoder, error) {
	if dec, ok := supportedFormats[format]; ok {
		return dec, nil
	}

	return nil, ErrUnsupported
}

func (ad *AudioFile) Duration() time.Duration {
	return ad.buffer.Format().SampleRate.D(ad.buffer.Len())
}

// TODO: Implement this
func (ad *AudioFile) Play()   {}
func (ad *AudioFile) Pause()  {}
func (ad *AudioFile) Resume() {}
