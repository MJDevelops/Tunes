package audio

import (
	"context"
	"errors"
	"io"
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
	Metadata TrackMeta
	buffer   *beep.Buffer
	format   beep.Format
}

type PlayingQueue struct {
	Queue []AudioFile
	ctx   context.Context
}

var (
	ErrUnsupported = errors.New("unsupported file format")
)

func NewAudioFile(path string) (AudioFile, error) {
	af := AudioFile{}

	var (
		err      error
		format   beep.Format
		buffer   *beep.Buffer
		streamer beep.StreamSeekCloser
	)

	f, err := os.Open(path)
	if err != nil {
		return af, err
	}

	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".flac":
		streamer, format, err = flac.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	default:
		return af, ErrUnsupported
	}

	if err != nil {
		return af, err
	}

	buffer = beep.NewBuffer(format)
	buffer.Append(streamer)

	f.Seek(0, io.SeekStart)
	if ext == ".wav" {
		af.Metadata, _ = parseWavMeta(f)
	} else {
		af.Metadata, _ = parseTagMeta(f, buffer)
	}

	streamer.Close()

	af.buffer = buffer
	af.format = format
	af.Path = path

	return af, nil
}

func (pq *PlayingQueue) SetContext(ctx context.Context) {
	pq.ctx = ctx
}
