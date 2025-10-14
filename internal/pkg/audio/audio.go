package audio

import (
	"errors"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
)

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

func NewAudioFile(decoder Decoder) (AudioFile, error) {
	af := AudioFile{}

	var (
		err      error
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
