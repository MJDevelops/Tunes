package audio

import (
	"os"
	"time"

	"github.com/go-audio/wav"
	"github.com/gopxl/beep/v2"
	beepwav "github.com/gopxl/beep/v2/wav"
	tunesos "github.com/mjdevelops/tunes/internal/pkg/os"
)

type WavDecoder struct {
	path     string
	duration time.Duration
}

func (wd *WavDecoder) New(path string) (Decoder, error) {
	if tunesos.GetFileExtension(path) != ".wav" {
		return nil, ErrUnsupported
	}

	return &WavDecoder{
		path: path,
	}, nil
}

func (wd *WavDecoder) Decode() (beep.StreamSeekCloser, beep.Format, error) {
	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
		err      error
	)

	if tunesos.GetFileExtension(wd.path) != ".wav" {
		return nil, format, ErrUnsupported
	}

	file, err := os.Open(wd.path)
	if err != nil {
		return nil, format, err
	}
	defer file.Close()

	streamer, format, err = beepwav.Decode(file)
	if err != nil {
		return nil, format, err
	}

	wd.duration = format.SampleRate.D(streamer.Len())

	return streamer, format, nil
}

func (wd *WavDecoder) ParseMeta() (TrackMeta, error) {
	if tunesos.GetFileExtension(wd.path) != ".wav" {
		return TrackMeta{}, ErrUnsupported
	}

	file, err := os.Open(wd.path)
	if err != nil {
		return TrackMeta{}, err
	}
	defer file.Close()

	dec := wav.NewDecoder(file)
	dec.ReadMetadata()

	return TrackMeta{
		Title:  dec.Metadata.Title,
		Album:  dec.Metadata.Product,
		Genre:  dec.Metadata.Genre,
		Artist: dec.Metadata.Artist,
	}, nil
}

func (wd *WavDecoder) Duration() time.Duration {
	return wd.duration
}
