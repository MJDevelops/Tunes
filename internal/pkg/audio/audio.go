package audio

import (
	"errors"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
	"github.com/gopxl/beep/v2/speaker"
)

type Decoder interface {
	New(path string) (Decoder, error)
	DecodeAudio() (*AudioFile, error)
	ParseMeta() (TrackMeta, error)
}

type AudioFile struct {
	buffer   *beep.Buffer
	format   beep.Format
	vol      *effects.Volume
	ctrl     *beep.Ctrl
	streamer beep.StreamSeeker
	stop     chan struct{}
}

var supportedFormats = make(map[string]Decoder)

var (
	ErrUnsupported = errors.New("unsupported file format")
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

func (ad *AudioFile) Duration() time.Duration {
	return ad.buffer.Format().SampleRate.D(ad.buffer.Len())
}

func (ad *AudioFile) Play(volume float64) (pos <-chan int) {
	position := make(chan int)
	done := make(chan struct{})

	ad.stop = make(chan struct{})
	ad.streamer = ad.buffer.Streamer(0, ad.buffer.Len())
	ad.ctrl = &beep.Ctrl{Streamer: ad.streamer, Paused: false}
	ad.vol = &effects.Volume{
		Streamer: ad.ctrl,
		Base:     2,
		Volume:   volume,
		Silent:   false,
	}

	speaker.Init(ad.format.SampleRate, ad.format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(ad.vol, beep.Callback(func() {
		done <- struct{}{}
	})))

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				speaker.Lock()
				position <- int(ad.format.SampleRate.D(ad.streamer.Position()).Seconds())
				speaker.Unlock()
			case <-ad.stop:
				speaker.Lock()
				speaker.Clear()
				speaker.Suspend()
				ad.vol, ad.ctrl, ad.stop, ad.streamer = nil, nil, nil, nil
				speaker.Unlock()
				return
			case <-done:
				ad.TogglePlayback()
				speaker.Lock()
				ad.streamer.Seek(0)
				speaker.Unlock()
			}
		}
	}()

	return position
}

func (ad *AudioFile) TogglePlayback() {
	speaker.Lock()
	defer speaker.Unlock()
	if ad.ctrl != nil {
		ad.ctrl.Paused = !ad.ctrl.Paused
	}
}

func (ad *AudioFile) Seek(d time.Duration) {
	speaker.Lock()
	defer speaker.Unlock()
	if ad.streamer != nil {
		ad.streamer.Seek(ad.format.SampleRate.N(d))
	}
}

func (ad *AudioFile) Stop() {
	if ad.stop != nil {
		ad.stop <- struct{}{}
	}
}

func (ad *AudioFile) Volume(vol float64) {
	speaker.Lock()
	defer speaker.Unlock()
	if ad.vol != nil {
		ad.vol.Volume = vol
	}
}
