package sound

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

type AudioFile struct {
	Path     string
	Format   beep.Format
	streamer beep.StreamSeekCloser
}

type PlayingQueue struct {
	Queue []AudioFile
	ctx   context.Context
}

var (
	ErrUnsupported = errors.New("unsupported file format")
)

func NewAudioFile(path string) (*AudioFile, error) {
	af := &AudioFile{}
	var err error
	var format beep.Format
	var streamer beep.StreamSeekCloser

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	default:
		return nil, ErrUnsupported
	}

	if err != nil {
		return nil, err
	}

	af.streamer = streamer
	af.Format = format
	af.Path = path

	return af, nil
}

func (pq *PlayingQueue) SetContext(ctx context.Context) {
	pq.ctx = ctx
}
