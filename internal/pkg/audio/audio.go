package audio

// #cgo CFLAGS: -g -Wall
// #cgo pkg-config: libavformat libavcodec libavutil libswresample
import "C"
import (
	"errors"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/speaker"

	"github.com/mjdevelops/tunes/internal/pkg/os"
)

type Decoder interface {
	New(path string) (Decoder, error)
	Decode() (beep.StreamSeekCloser, beep.Format, error)
	ParseMeta() (TrackMeta, error)
	Duration() time.Duration
}

type AudioSink struct {
	format   beep.Format
	vol      *effects.Volume
	ctrl     *beep.Ctrl
	streamer beep.StreamSeeker
	decoder  Decoder
	stop     chan struct{}
	done     bool
}

var supportedFormats = make(map[string]Decoder)

var (
	ErrUnsupported = errors.New("unsupported file format")
	ErrPlaying     = errors.New("sink is still playing audio")
)

func RegisterDecoder(decoder Decoder, formats ...string) {
	for _, format := range formats {
		supportedFormats[format] = decoder
	}
}

func GetDecoder(format string) (Decoder, error) {
	if dec, ok := supportedFormats[format]; ok {
		return dec, nil
	}

	return nil, ErrUnsupported
}

func NewAudioSink() *AudioSink {
	return &AudioSink{
		done: true,
		stop: make(chan struct{}),
	}
}

func (a *AudioSink) Duration() time.Duration {
	return a.decoder.Duration()
}

func (a *AudioSink) Init(trackPath string) error {
	if !a.done {
		return ErrPlaying
	}

	decoder, err := GetDecoder(os.GetFileExtension(trackPath))
	if err != nil {
		// Use libav as fallback
		decoder = &AVDecoder{}
	}

	a.streamer, a.format, err = decoder.Decode()
	if err != nil {
		return err
	}

	a.decoder = decoder

	a.ctrl = &beep.Ctrl{Streamer: a.streamer, Paused: false}
	a.vol = &effects.Volume{
		Streamer: a.ctrl,
		Base:     2,
		Volume:   0,
		Silent:   true,
	}

	speaker.Init(a.format.SampleRate, a.format.SampleRate.N(time.Second/10))

	return nil
}

func (a *AudioSink) Play(volume float64) (pos <-chan int) {
	if !a.done {
		return nil
	}

	position := make(chan int)
	done := make(chan struct{})

	a.vol.Volume = volume
	a.vol.Silent = false
	a.done = false

	speaker.Play(beep.Seq(a.vol, beep.Callback(func() {
		a.done = true
		done <- struct{}{}
	})))

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				speaker.Lock()
				position <- int(a.format.SampleRate.D(a.streamer.Position()).Seconds())
				speaker.Unlock()
			case <-a.stop:
				speaker.Lock()
				speaker.Clear()
				speaker.Suspend()
				speaker.Unlock()
				return
			case <-done:
				a.TogglePlayback()
				speaker.Lock()
				a.streamer.Seek(0)
				speaker.Unlock()
			}
		}
	}()

	return position
}

func (a *AudioSink) TogglePlayback() {
	speaker.Lock()
	defer speaker.Unlock()
	if a.ctrl != nil {
		a.ctrl.Paused = !a.ctrl.Paused
	}
}

func (a *AudioSink) Seek(d time.Duration) {
	speaker.Lock()
	defer speaker.Unlock()
	if a.streamer != nil {
		a.streamer.Seek(a.format.SampleRate.N(d))
	}
}

func (a *AudioSink) Stop() {
	if a.stop != nil {
		a.stop <- struct{}{}
	}
}

func (a *AudioSink) Volume(vol float64) {
	speaker.Lock()
	defer speaker.Unlock()
	if a.vol != nil {
		a.vol.Volume = vol
	}
}
