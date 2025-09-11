package audio

import (
	"os"

	"github.com/go-audio/wav"
)

func parseWavMeta(file *os.File) (TrackMeta, error) {
	dec := wav.NewDecoder(file)
	dec.ReadMetadata()
	dur, _ := dec.Duration()
	dec.Rewind()

	return TrackMeta{
		Title:    dec.Metadata.Title,
		Album:    dec.Metadata.Product,
		Duration: dur,
		Genre:    dec.Metadata.Genre,
		Artist:   dec.Metadata.Artist,
	}, nil
}
