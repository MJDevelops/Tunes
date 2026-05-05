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

type Resampler struct {
	src            *astiav.SoftwareResampleContext
	formatContext  *astiav.FormatContext
	pkt            *astiav.Packet
	frame          *astiav.Frame
	resampledFrame *astiav.Frame
	finalFrame     *astiav.Frame
	af             *astiav.AudioFifo
}

func (r *Resampler) ResampleStereo(file string) (*Audio, error) {
	audio := &Audio{}

	r.src = astiav.AllocSoftwareResampleContext()
	defer r.src.Free()

	r.formatContext = astiav.AllocFormatContext()
	defer r.formatContext.Free()

	r.pkt = astiav.AllocPacket()
	defer r.pkt.Free()

	r.frame = astiav.AllocFrame()
	defer r.frame.Free()

	r.resampledFrame = astiav.AllocFrame()
	defer r.resampledFrame.Free()

	r.resampledFrame.SetChannelLayout(astiav.ChannelLayoutStereo)
	r.resampledFrame.SetSampleFormat(astiav.SampleFormatDblp)
	r.resampledFrame.SetSampleRate(44100)
	r.resampledFrame.SetNbSamples(r.resampledFrame.SampleRate() / 25)

	if err := r.resampledFrame.AllocBuffer(0); err != nil {
		return nil, fmt.Errorf("allocating resampled frame buffer failed: %w", err)
	}

	r.finalFrame = astiav.AllocFrame()
	defer r.finalFrame.Free()

	r.finalFrame.SetChannelLayout(r.resampledFrame.ChannelLayout())
	r.finalFrame.SetNbSamples(r.resampledFrame.NbSamples())
	r.finalFrame.SetSampleFormat(r.resampledFrame.SampleFormat())
	r.finalFrame.SetSampleRate(r.resampledFrame.SampleRate())
	if err := r.finalFrame.AllocBuffer(0); err != nil {
		return nil, fmt.Errorf("allocating final frame buffer failed: %w", err)
	}

	r.af = astiav.AllocAudioFifo(r.finalFrame.SampleFormat(), r.finalFrame.ChannelLayout().Channels(), r.finalFrame.NbSamples())
	defer r.af.Free()

	r.formatContext.OpenInput(file, nil, nil)
	defer r.formatContext.CloseInput()

	r.formatContext.FindStreamInfo(nil)

	for _, stream := range r.formatContext.Streams() {
		codecParameters := stream.CodecParameters()
		if codecParameters.MediaType() != astiav.MediaTypeAudio {
			continue
		}
		codec := astiav.FindDecoder(codecParameters.CodecID())
		codecContext := astiav.AllocCodecContext(codec)
		codecParameters.ToCodecContext(codecContext)

		for {
			if stop := func() bool {
				if err := r.formatContext.ReadFrame(r.pkt); err != nil {
					if !errors.Is(err, astiav.ErrEof) {
						log.Println(fmt.Errorf("reading frame failed: %w", err))
					}
					return true
				}

				defer r.pkt.Unref()

				if r.pkt.StreamIndex() != stream.Index() {
					return true
				}

				if err := codecContext.SendPacket(r.pkt); err != nil {
					if !errors.Is(err, astiav.ErrEof) && !errors.Is(err, astiav.ErrEagain) {
						log.Println(fmt.Errorf("sending packet failed: %w", err))
					}
					return true
				}

				for {
					if stop := func() bool {
						if err := codecContext.ReceiveFrame(r.frame); err != nil {
							if !errors.Is(err, astiav.ErrEof) && !errors.Is(err, astiav.ErrEagain) {
								log.Println(fmt.Errorf("receiving frame failed: %w", err))
							}
							return true
						}

						defer r.frame.Unref()

						if err := r.src.ConvertFrame(r.frame, r.resampledFrame); err != nil {
							log.Println(fmt.Errorf("resampling frame failed: %w", err))
							return true
						}

						if nbSamples := r.resampledFrame.NbSamples(); nbSamples > 0 {
							if err := r.addResampledFrameToFifo(false); err != nil {
								log.Println(fmt.Errorf("adding resampled frame to fifo failed: %w", err))
								return true
							}

							if err := r.flushSoftwareResampleContext(false); err != nil {
								log.Println(fmt.Errorf("flushing software resample context failed: %w", err))
								return true
							}
						}
						return false
					}(); stop {
						break
					}
				}

				return false
			}(); stop {
				break
			}
		}
	}

	if err := r.flushSoftwareResampleContext(true); err != nil {
		return nil, fmt.Errorf("flushing software resample context failed: %w", err)
	}

	return audio, nil
}

func (r *Resampler) addResampledFrameToFifo(flush bool) error {
	if r.resampledFrame.NbSamples() > 0 {
		if _, err := r.af.Write(r.resampledFrame); err != nil {
			return fmt.Errorf("writing to audio fifo failed: %w", err)
		}
	}

	for {
		if (flush && r.af.Size() > 0) || (!flush && r.af.Size() >= r.finalFrame.NbSamples()) {
			n, err := r.af.Read(r.finalFrame)
			if err != nil {
				return fmt.Errorf("reading from audio fifo failed: %w", err)
			}
			r.finalFrame.SetNbSamples(n)
			continue
		}
		break
	}
	return nil
}

func (r *Resampler) flushSoftwareResampleContext(finalFlush bool) error {
	for {
		if finalFlush || r.src.Delay(int64(r.resampledFrame.SampleRate())) >= int64(r.resampledFrame.NbSamples()) {
			if err := r.src.ConvertFrame(nil, r.resampledFrame); err != nil {
				return fmt.Errorf("flushing software resample context failed: %w", err)
			}

			if r.resampledFrame.NbSamples() == 0 {
				break
			}

			if err := r.addResampledFrameToFifo(finalFlush); err != nil {
				return fmt.Errorf("adding resampled frame to fifo failed: %w", err)
			}
			continue
		}
		break
	}
	return nil
}
