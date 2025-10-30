package audio

import (
	"os"

	"github.com/go-audio/wav"
	"github.com/gopxl/beep/v2"
	beepwav "github.com/gopxl/beep/v2/wav"
	tunesos "github.com/mjdevelops/tunes/internal/pkg/os"
)

type WavDecoder struct {
	path string
}

func (wd *WavDecoder) New(path string) (Decoder, error) {
	if tunesos.GetFileExtension(path) != ".wav" {
		return nil, ErrUnsupported
	}

	return &WavDecoder{
		path: path,
	}, nil
}

func (wd *WavDecoder) DecodeAudio() (*AudioFile, error) {
	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
		buffer   *beep.Buffer
		err      error
	)

	if tunesos.GetFileExtension(wd.path) != ".wav" {
		return nil, ErrUnsupported
	}

	file, err := os.Open(wd.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	streamer, format, err = beepwav.Decode(file)
	if err != nil {
		return nil, err
	}

	buffer = beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	return &AudioFile{
		buffer: buffer,
		format: format,
	}, nil
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
