package audio

import (
	"io"
	"os"

	"github.com/go-audio/wav"
	"github.com/gopxl/beep/v2"
	beepwav "github.com/gopxl/beep/v2/wav"
	"github.com/mjdevelops/tunes/internal/pkg/util"
)

type WavDecoder struct {
	file *os.File
}

func NewWavDecoder(path string) (*WavDecoder, error) {
	var err error

	wd := &WavDecoder{}
	if util.GetFileExtension(path) != ".wav" {
		return nil, ErrUnsupported
	}

	wd.file, err = os.Open(path)
	return wd, err
}

func (wd *WavDecoder) DecodeAudio() (beep.StreamSeekCloser, beep.Format, error) {
	return beepwav.Decode(wd.file)
}

func (wd *WavDecoder) ParseMeta() TrackMeta {
	wd.file.Seek(0, io.SeekStart)

	dec := wav.NewDecoder(wd.file)
	dec.ReadMetadata()
	dec.Rewind()

	return TrackMeta{
		Title:  dec.Metadata.Title,
		Album:  dec.Metadata.Product,
		Genre:  dec.Metadata.Genre,
		Artist: dec.Metadata.Artist,
	}
}
