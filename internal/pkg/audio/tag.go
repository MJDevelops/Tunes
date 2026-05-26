package audio

import (
	"os"
	"slices"
	"time"

	"github.com/dhowden/tag"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"
	tunesos "github.com/mjdevelops/tunes/internal/pkg/os"
)

// TagDecoder represents an audio file decoder for formats supported by github.com/dhowden/tag
type TagDecoder struct {
	path     string
	duration time.Duration
}

var tagFormats = []string{".flac", ".ogg", ".mp3"}

func (td *TagDecoder) Init(path string) error {
	fileExt := tunesos.GetFileExtension(path)

	if !slices.Contains(tagFormats, fileExt) {
		return ErrUnsupported
	}

	td.path = path
	return nil
}

func (td *TagDecoder) Decode() (beep.StreamSeekCloser, beep.Format, error) {
	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
		err      error
	)
	fileExt := tunesos.GetFileExtension(td.path)

	if !slices.Contains(tagFormats, fileExt) {
		return nil, format, ErrUnsupported
	}

	file, err := os.Open(td.path)
	if err != nil {
		return nil, format, err
	}
	defer file.Close()

	switch fileExt {
	case ".flac":
		streamer, format, err = flac.Decode(file)
	case ".ogg":
		streamer, format, err = vorbis.Decode(file)
	case ".mp3":
		streamer, format, err = mp3.Decode(file)
	}

	if err != nil {
		return nil, format, err
	}

	td.duration = format.SampleRate.D(streamer.Len())

	return streamer, format, nil
}

func (td *TagDecoder) ParseMeta() (TrackMeta, error) {
	trackMeta := TrackMeta{}

	file, err := os.Open(td.path)
	if err != nil {
		return TrackMeta{}, err
	}
	defer file.Close()

	meta, err := tag.ReadFrom(file)
	if err != nil {
		return trackMeta, err
	}

	trackMeta.Title = meta.Title()
	trackMeta.Album = meta.Album()
	trackMeta.Artist = meta.Artist()
	trackMeta.Genre = meta.Genre()

	return trackMeta, nil
}

func (td *TagDecoder) Duration() time.Duration {
	return td.duration
}
