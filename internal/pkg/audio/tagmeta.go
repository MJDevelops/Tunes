package audio

import (
	"os"

	"github.com/dhowden/tag"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"
	tunesos "github.com/mjdevelops/tunes/internal/pkg/os"
)

// TagDecoder represents an audio file decoder for formats supported by github.com/dhowden/tag
type TagDecoder struct {
	file    *os.File
	fileExt string
}

func NewTagDecoder(path string) (*TagDecoder, error) {
	td := &TagDecoder{}
	td.fileExt = tunesos.GetFileExtension(path)

	if err := IsSupportedFormat(td.fileExt); err != nil {
		return nil, err

	}

	td.file, _ = os.Open(path)
	return td, nil
}

func (td *TagDecoder) DecodeAudio() (beep.StreamSeekCloser, beep.Format, error) {
	switch td.fileExt {
	case ".flac":
		return flac.Decode(td.file)
	case ".ogg":
		return vorbis.Decode(td.file)
	case ".mp3":
		return mp3.Decode(td.file)
	}

	return nil, beep.Format{}, ErrUnsupported
}

func (td *TagDecoder) ParseMeta() (TrackMeta, error) {
	trackMeta := TrackMeta{}
	meta, err := tag.ReadFrom(td.file)
	if err != nil {
		return trackMeta, err
	}

	trackMeta.Title = meta.Title()
	trackMeta.Album = meta.Album()
	trackMeta.Artist = meta.Artist()
	trackMeta.Genre = meta.Genre()

	return trackMeta, nil
}
