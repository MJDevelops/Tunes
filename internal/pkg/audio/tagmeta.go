package audio

import (
	"os"
	"time"

	"github.com/dhowden/tag"
	"github.com/gopxl/beep/v2"
)

// Represents Metadata supported by dhowden/tag. (OGG, FLAC, MP3, MP4)
type TagMeta struct {
	title    string
	duration time.Duration
	artist   string
	album    string
	genre    string
}

func (f *TagMeta) Album() string {
	return f.album
}

func (f *TagMeta) Artist() string {
	return f.artist
}

func (f *TagMeta) Genre() string {
	return f.genre
}

func (f *TagMeta) Duration() time.Duration {
	return f.duration
}

func (f *TagMeta) Title() string {
	return f.title
}

func parseTagMeta(file *os.File, buf *beep.Buffer, format beep.Format) (TrackMeta, error) {
	meta, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	duration := format.SampleRate.D(buf.Len())

	return &TagMeta{
		title:    meta.Title(),
		album:    meta.Album(),
		artist:   meta.Artist(),
		genre:    meta.Genre(),
		duration: duration,
	}, nil
}
