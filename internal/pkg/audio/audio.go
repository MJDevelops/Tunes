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
	Init(path string) error
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
	stopped  bool
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

func GetDecoder(file string) (Decoder, error) {
	format := os.GetFileExtension(file)
	if dec, ok := supportedFormats[format]; ok {
		return dec, nil
	}

	return &AVDecoder{}, nil
}

func NewAudioSink() *AudioSink {
	return &AudioSink{
		stopped: true,
		stop:    make(chan struct{}, 1),
	}
}

func (a *AudioSink) Duration() time.Duration {
	return a.decoder.Duration()
}

func (a *AudioSink) Init(trackPath string) error {
	var err error

	if a.IsPlaying() {
		a.Stop()
	}

	decoder, err := GetDecoder(trackPath)
	if err != nil {
		return err
	}

	err = decoder.Init(trackPath)
	if err != nil {
		return err
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
		Silent:   false,
	}

	speaker.Init(a.format.SampleRate, a.format.SampleRate.N(time.Second/10))

	return nil
}

func (a *AudioSink) Play(volume float64) error {
	if a.IsPlaying() {
		return ErrPlaying
	}

	a.vol.Volume = volume
	a.stopped = false

	speaker.Play(beep.Seq(a.vol, beep.Callback(func() {
		a.stop <- struct{}{}
	})))

	go func() {
		<-a.stop
		a.stopPlayback()
		a.Seek(0)
	}()

	return nil
}

func (a *AudioSink) SetPlayback(playing bool) {
	speaker.Lock()
	defer speaker.Unlock()
	if a.ctrl != nil {
		a.ctrl.Paused = !playing
	}
}

func (a *AudioSink) IsPaused() bool {
	speaker.Lock()
	defer speaker.Unlock()
	if a.ctrl != nil {
		return a.ctrl.Paused
	}

	return false
}

func (a *AudioSink) Seek(d time.Duration) {
	speaker.Lock()
	defer speaker.Unlock()
	if a.streamer != nil {
		a.streamer.Seek(a.format.SampleRate.N(d))
	}
}

func (a *AudioSink) Stop() {
	speaker.Lock()
	defer speaker.Unlock()
	if a.stop != nil && len(a.stop) == 0 && a.stopped == false {
		a.stop <- struct{}{}
	}
}

func (a *AudioSink) SetVolume(vol float64) {
	speaker.Lock()
	defer speaker.Unlock()
	if a.vol != nil {
		a.vol.Volume = vol
	}
}

func (a *AudioSink) Position() time.Duration {
	speaker.Lock()
	defer speaker.Unlock()
	if a.streamer != nil {
		return a.format.SampleRate.D(a.streamer.Position())
	}
	return 0
}

func (a *AudioSink) IsPlaying() bool {
	speaker.Lock()
	defer speaker.Unlock()
	return !a.stopped
}

func (a *AudioSink) Volume() float64 {
	speaker.Lock()
	defer speaker.Unlock()
	if a.vol != nil {
		return a.vol.Volume
	}

	return 0
}

func (a *AudioSink) SetMute(muted bool) {
	speaker.Lock()
	defer speaker.Unlock()
	if a.vol != nil {
		a.vol.Silent = muted
	}
}

func (a *AudioSink) stopPlayback() {
	speaker.Lock()
	defer speaker.Unlock()
	speaker.Clear()
	a.stopped = true
}
