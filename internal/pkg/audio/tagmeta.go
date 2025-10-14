package audio

import (
	"os"
	"slices"

	"github.com/dhowden/tag"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/mjdevelops/tunes/internal/pkg/util"
)

// Audio file decoder for formats supported by github.com/dhowden/tag
type TagDecoder struct {
	file    *os.File
	fileExt string
}

var TagFormats = []string{".flac", ".ogg", ".mp3"}

func NewTagDecoder(path string) (*TagDecoder, error) {
	td := &TagDecoder{}
	td.fileExt = util.GetFileExtension(path)

	if slices.Contains(TagFormats, td.fileExt) {
		td.file, _ = os.Open(path)
		return td, nil
	}

	return nil, ErrUnsupported
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
