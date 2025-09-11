package audio

import (
	"os"

	"github.com/dhowden/tag"
	"github.com/gopxl/beep/v2"
)

func parseTagMeta(file *os.File, buf *beep.Buffer) (TrackMeta, error) {
	trackMeta := TrackMeta{}
	meta, err := tag.ReadFrom(file)
	if err != nil {
		return trackMeta, err
	}
	duration := buf.Format().SampleRate.D(buf.Len())

	trackMeta.Title = meta.Title()
	trackMeta.Album = meta.Album()
	trackMeta.Artist = meta.Artist()
	trackMeta.Genre = meta.Genre()
	trackMeta.Duration = duration

	return trackMeta, nil
}
