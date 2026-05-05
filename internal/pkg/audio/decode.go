package audio

import (
	"errors"
	"fmt"
	"log"

	"github.com/asticode/go-astiav"
	"github.com/asticode/go-astikit"
	"github.com/gopxl/beep/v2"
)

type FFDecoder struct {
	src            *astiav.SoftwareResampleContext
	formatContext  *astiav.FormatContext
	pkt            *astiav.Packet
	frame          *astiav.Frame
	resampledFrame *astiav.Frame
	finalFrame     *astiav.Frame
	af             *astiav.AudioFifo
	bytesPerFrame  int
	format         beep.Format
	ad             *astiav.Stream
	codec          *astiav.Codec
	codecContext   *astiav.CodecContext
	pos            int
	err            error
	buf            []byte
	decLastPTS     *int64
}

func Decode(path string) (s beep.StreamSeekCloser, format beep.Format, err error) {
	d, err := NewDecoder(path)
	if err != nil {
		return nil, beep.Format{}, fmt.Errorf("creating decoder failed: %w", err)
	}

	format = beep.Format{
		SampleRate:  beep.SampleRate(d.SampleRate()),
		NumChannels: d.Channels(),
		Precision:   d.Precision(),
	}

	d.format = format

	return d, format, nil
}

func NewDecoder(path string) (*FFDecoder, error) {
	r := &FFDecoder{}
	r.src = astiav.AllocSoftwareResampleContext()
	r.formatContext = astiav.AllocFormatContext()
	r.pkt = astiav.AllocPacket()
	r.frame = astiav.AllocFrame()
	r.resampledFrame = astiav.AllocFrame()
	r.finalFrame = astiav.AllocFrame()

	r.resampledFrame.SetChannelLayout(astiav.ChannelLayoutStereo)
	r.resampledFrame.SetSampleFormat(astiav.SampleFormatS16P)
	r.resampledFrame.SetSampleRate(44100)
	r.resampledFrame.SetNbSamples(r.resampledFrame.SampleRate() / 25)

	if err := r.resampledFrame.AllocBuffer(0); err != nil {
		return nil, fmt.Errorf("allocating resampled frame buffer failed: %w", err)
	}

	r.finalFrame.SetChannelLayout(r.resampledFrame.ChannelLayout())
	r.finalFrame.SetNbSamples(r.resampledFrame.NbSamples())
	r.finalFrame.SetSampleFormat(r.resampledFrame.SampleFormat())
	r.finalFrame.SetSampleRate(r.resampledFrame.SampleRate())

	if err := r.finalFrame.AllocBuffer(0); err != nil {
		return nil, fmt.Errorf("allocating final frame buffer failed: %w", err)
	}

	r.bytesPerFrame = r.finalFrame.ChannelLayout().Channels() * r.finalFrame.SampleFormat().BytesPerSample()

	r.af = astiav.AllocAudioFifo(r.finalFrame.SampleFormat(), r.finalFrame.ChannelLayout().Channels(), r.finalFrame.NbSamples())

	r.formatContext.OpenInput(path, nil, nil)
	r.formatContext.FindStreamInfo(nil)

	for _, stream := range r.formatContext.Streams() {
		codecParameters := stream.CodecParameters()
		if codecParameters.MediaType() != astiav.MediaTypeAudio {
			continue
		}
		r.ad = stream
		break
	}

	codecParameters := r.ad.CodecParameters()
	r.codec = astiav.FindDecoder(codecParameters.CodecID())
	r.codecContext = astiav.AllocCodecContext(r.codec)
	codecParameters.ToCodecContext(r.codecContext)
	r.codecContext.Open(r.codec, nil)

	return r, nil
}

func (r *FFDecoder) Err() error {
	return r.err
}

func (r *FFDecoder) Close() error {
	r.src.Free()
	r.formatContext.Free()
	r.pkt.Free()
	r.frame.Free()
	r.resampledFrame.Free()
	r.finalFrame.Free()
	r.af.Free()
	r.formatContext.CloseInput()
	r.codecContext.Free()

	return nil
}

func (r *FFDecoder) SampleRate() int {
	return r.finalFrame.SampleRate()
}

func (r *FFDecoder) Channels() int {
	return r.finalFrame.ChannelLayout().Channels()
}

func (r *FFDecoder) Precision() int {
	return r.finalFrame.SampleFormat().BytesPerSample()
}

func (r *FFDecoder) Position() int {
	return r.pos / r.bytesPerFrame
}

func (r *FFDecoder) Length() int64 {
	return r.ad.NbFrames() * int64(r.bytesPerFrame)
}

func (r *FFDecoder) Len() int {
	return int(r.Length()) / r.bytesPerFrame
}

// TODO
func (r *FFDecoder) Stream(samples [][2]float64) (n int, ok bool) {
	buf := make([]byte, r.bytesPerFrame)

	for {
		if stop := func() bool {
			if len(r.buf) > 0 {
				return true
			}

			if err := r.formatContext.ReadFrame(r.pkt); err != nil {
				if !errors.Is(err, astiav.ErrEof) {
					log.Println(fmt.Errorf("reading frame failed: %w", err))
				}
				return true
			}

			defer r.pkt.Unref()

			if r.pkt.StreamIndex() != r.ad.Index() {
				return false
			}

			if err := r.codecContext.SendPacket(r.pkt); err != nil {
				if !errors.Is(err, astiav.ErrEof) && !errors.Is(err, astiav.ErrEagain) {
					log.Println(fmt.Errorf("sending packet failed: %w", err))
				}
				return true
			}

			for {
				if stop := func() bool {
					if err := r.codecContext.ReceiveFrame(r.frame); err != nil {
						if !errors.Is(err, astiav.ErrEof) && !errors.Is(err, astiav.ErrEagain) {
							log.Println(fmt.Errorf("receiving frame failed: %w", err))
						}
						return true
					}

					defer r.frame.Unref()

					if r.decLastPTS != nil && *r.decLastPTS >= r.frame.Pts() {
						return false
					}
					r.decLastPTS = astikit.Int64Ptr(r.frame.Pts())

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

						r.copyFrameDataToBuffer()
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

	if err := r.flushSoftwareResampleContext(true); err != nil {
		r.err = fmt.Errorf("flushing software resample context failed: %w", err)
	}

	if r.finalFrame.NbSamples() > 0 {
		r.copyFrameDataToBuffer()
	}

	for i := range samples {
		b := copy(buf, r.buf)
		r.buf = r.buf[b:]
		if b == len(buf) {
			samples[i], _ = r.format.DecodeSigned(buf[:])
			n++
			r.pos += b
			ok = true
		}
	}

	return n, ok
}

// TODO
func (r *FFDecoder) Seek(p int) error {
	return nil
}

func (r *FFDecoder) copyFrameDataToBuffer() error {
	var err error
	var fBuf []byte
	if fBuf, err = r.finalFrame.Data().Bytes(0); err != nil {
		return fmt.Errorf("copying samples to buffer failed: %w", err)
	}

	r.buf = append(r.buf, fBuf...)
	return nil
}

func (r *FFDecoder) addResampledFrameToFifo(flush bool) error {
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

func (r *FFDecoder) flushSoftwareResampleContext(finalFlush bool) error {
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
