package audio

import (
	"container/list"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/mjdevelops/tunes/internal/pkg/os"
)

var supportedFormats = []string{".flac", ".ogg", ".mp3", ".wav"}

type Decoder interface {
	DecodeAudio() (beep.StreamSeekCloser, beep.Format, error)
	ParseMeta() (TrackMeta, error)
}

type AudioFile struct {
	Metadata TrackMeta
	buffer   *beep.Buffer
	format   beep.Format
}

type Queue struct {
	list    *list.List
	current *list.Element
	mu      sync.Mutex
}

var (
	ErrUnsupported = errors.New("unsupported file format")
)

// NewAudioFile constructs a new AudioFile struct with the provided decoder.
//
// This function returns err != nil if the audio decoding fails.
func NewAudioFile(decoder Decoder) (ad AudioFile, err error) {
	af := AudioFile{}

	var (
		format   beep.Format
		buffer   *beep.Buffer
		streamer beep.StreamSeekCloser
	)

	streamer, format, err = decoder.DecodeAudio()

	if err != nil {
		return af, err
	}

	buffer = beep.NewBuffer(format)
	buffer.Append(streamer)
	af.Metadata, _ = decoder.ParseMeta()
	streamer.Close()

	af.buffer = buffer
	af.format = format

	return af, nil
}

func NewDecoder(file string) (Decoder, error) {
	ext := os.GetFileExtension(file)
	if err := IsSupportedFormat(ext); err != nil {
		return nil, err
	}

	switch ext {
	case ".wav":
		return NewWavDecoder(file)
	default:
		return NewTagDecoder(file)
	}
}

func (ad *AudioFile) Duration() time.Duration {
	return ad.buffer.Format().SampleRate.D(ad.buffer.Len())
}

// TODO: Implement this
func (ad *AudioFile) Play()  {}
func (ad *AudioFile) Pause() {}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

func (q *Queue) Add(ad *AudioFile) *list.Element {
	q.mu.Lock()
	defer q.mu.Unlock()

	e := q.list.PushBack(ad)

	if q.current == nil {
		q.current = e
	}

	return e
}

func (q *Queue) AddAfter(ad *AudioFile, e *list.Element) *list.Element {
	return q.list.InsertAfter(ad, e)
}

func (q *Queue) Next() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.current != nil {
		next := q.current.Next()

		if next != nil {
			q.current = next
		}
	}
}

func (q *Queue) Previous() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.current != nil {
		prev := q.current.Prev()

		if prev != nil {
			q.current = prev
		}

		q.current.Value.(*AudioFile).Play()
	}
}

func (q *Queue) MoveAfter(e *list.Element, mark *list.Element) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.list.MoveAfter(e, mark)
}

func (q *Queue) Play(e *list.Element) {
	q.mu.Lock()
	defer q.mu.Unlock()

	c := q.list.Front()
	if c == e {
		c.Value.(*AudioFile).Play()
	} else if q.current == e {
		q.current.Value.(*AudioFile).Play()
	} else {
		for {
			c = c.Next()
			if c == nil {
				break
			}

			if c == e {
				q.current = c
				c.Value.(*AudioFile).Play()
				break
			}
		}
	}
}

func (q *Queue) Pause() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.current != nil {
		q.current.Value.(*AudioFile).Pause()
	}
}

func (q *Queue) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.list.Init()
}

// IsSupportedFormat reports whether the provided format is supported
// in the scope of audio decoding. If the format is not supported this
// will return an error of type ErrUnsupported.
func IsSupportedFormat(format string) error {
	if !slices.Contains(supportedFormats, format) {
		return ErrUnsupported
	}
	return nil
}
