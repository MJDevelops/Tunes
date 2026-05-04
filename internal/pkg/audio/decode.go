package audio

import (
	"errors"
	"fmt"
	"log"

	"github.com/asticode/go-astiav"
)

type Audio struct {
	sampleRate int
	samples    [][2]float64
}

// TODO: implement this
func decodeAudio(file string) (*Audio, error) {
	audio := &Audio{}
	formatContext := astiav.AllocFormatContext()
	defer formatContext.Free()

	pkt := astiav.AllocPacket()
	defer pkt.Free()

	f := astiav.AllocFrame()
	defer f.Free()

	formatContext.OpenInput(file, nil, nil)
	defer formatContext.CloseInput()

	formatContext.FindStreamInfo(nil)

	for _, stream := range formatContext.Streams() {
		codecParameters := stream.CodecParameters()
		if codecParameters.MediaType() != astiav.MediaTypeAudio {
			continue
		}
		codec := astiav.FindDecoder(codecParameters.CodecID())
		codecContext := astiav.AllocCodecContext(codec)
		codecParameters.ToCodecContext(codecContext)

		for {
			if stop := func() bool {
				if err := formatContext.ReadFrame(pkt); err != nil {
					if !errors.Is(err, astiav.ErrEof) {
						log.Println(fmt.Errorf("reading frame failed: %w", err))
					}
					return true
				}

				defer pkt.Unref()

				if err := codecContext.SendPacket(pkt); err != nil {
					if !errors.Is(err, astiav.ErrEof) && !errors.Is(err, astiav.ErrEagain) {
						log.Println(fmt.Errorf("sending packet failed: %w", err))
					}
					return true
				}

				for {
					if stop := func() bool {
						if err := codecContext.ReceiveFrame(f); err != nil {
							if !errors.Is(err, astiav.ErrEof) && !errors.Is(err, astiav.ErrEagain) {
								log.Println(fmt.Errorf("receiving frame failed: %w", err))
							}
							return true
						}

						defer f.Unref()

						if f.SampleFormat() != astiav.SampleFormatDblp || f.ChannelLayout().Channels() != 2 {
							log.Printf("unsupported sample format: %s", f.SampleFormat().String())
							return true
						}

						if audio.sampleRate == 0 {
							audio.sampleRate = f.SampleRate()
						}

						for i := 0; i < f.NbSamples(); i++ {

						}

						return true
					}(); stop {
						break
					}
				}

				return true
			}(); stop {
				break
			}
		}
	}
	return audio, nil
}
