package audio

import (
	"os"
	"time"

	"github.com/go-audio/wav"
)

type WavMeta struct {
	title    string
	artist   string
	duration time.Duration
	album    string
	genre    string
}

func (w *WavMeta) Title() string {
	return w.title
}

func (w *WavMeta) Artist() string {
	return w.artist
}

func (w *WavMeta) Duration() time.Duration {
	return w.duration
}

func (w *WavMeta) Album() string {
	return w.album
}

func (w *WavMeta) Genre() string {
	return w.genre
}

func parseWavMeta(file *os.File) (*WavMeta, error) {
	dec := wav.NewDecoder(file)
	dec.ReadMetadata()
	dur, _ := dec.Duration()
	dec.Rewind()

	return &WavMeta{
		title:    dec.Metadata.Title,
		album:    dec.Metadata.Product,
		duration: dur,
		genre:    dec.Metadata.Genre,
		artist:   dec.Metadata.Artist,
	}, nil
}
