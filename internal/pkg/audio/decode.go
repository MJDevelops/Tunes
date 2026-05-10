package audio

// #include <stdlib.h>
// #include <stdint.h>
// #include <libavutil/avutil.h>
// #include "decoder.h"
// #include "sample_buffer.h"
import "C"
import (
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"

	"github.com/gopxl/beep/v2"
)

type AVDecoder struct {
	buf        *C.SampleBuffer
	dec        *C.Decoder
	f          beep.Format
	pos        int
	err        error
	samples    [2][]int16
	file       *C.char
	sampleRate int
}

const (
	avPrecision     = 2
	avChannels      = 2
	avBytesPerFrame = avChannels * avPrecision
)

var (
	ErrAlloc    = errors.New("error during allocation")
	ErrDecoding = errors.New("error during decoding")
	ErrSeek     = errors.New("error during seeking")
)

func NewDecoder(filename string) (s beep.StreamSeekCloser, format beep.Format, err error) {
	d := &AVDecoder{}
	d.file = C.CString(filename)
	d.dec = C.decoder_alloc(d.file)
	d.buf = C.sb_alloc()
	d.f = beep.Format{}

	if d.buf == nil || d.dec == nil {
		return nil, d.f, ErrAlloc
	}

	d.sampleRate = int(d.dec.sample_rate)

	d.f.Precision = avPrecision
	d.f.NumChannels = avChannels
	d.f.SampleRate = beep.SampleRate(d.sampleRate)

	return d, d.f, nil
}

func (d *AVDecoder) Samples() *[2][]int16 {
	return &d.samples
}

func (d *AVDecoder) SampleRate() int {
	return d.sampleRate
}

func (d *AVDecoder) NbSamples() int64 {
	if d.buf != nil {
		return int64(d.buf.channel_size)
	}
	return 0
}

func (d *AVDecoder) free() {
	if d.buf != nil {
		C.sb_free(&d.buf)
	}

	if d.dec != nil {
		C.decoder_free(&d.dec)
	}

	C.free(unsafe.Pointer(d.file))
}

func (d *AVDecoder) Stream(samples [][2]float64) (n int, ok bool) {
	b := make([]byte, avBytesPerFrame)
	for i := range samples {
		if len(d.samples[0]) == 0 {
			ret := C.decode(d.dec, d.buf, 1)
			if ret < 0 {
				return n, false
			}
			d.readBufToSamples(int(d.buf.channel_size))
		}

		// ffmpeg uses native endianness for audio samples
		binary.NativeEndian.PutUint16(b, uint16(d.samples[0][0]))
		binary.NativeEndian.PutUint16(b[2:], uint16(d.samples[1][0]))

		for i := range 2 {
			d.samples[i] = d.samples[i][1:]
		}

		samples[i], _ = d.f.DecodeSigned(b)
		d.pos += avBytesPerFrame
		n++
		ok = true
	}
	return n, ok
}

func (d *AVDecoder) readBufToSamples(numSamples int) {
	b := C.sb_flush(d.buf)
	for i := range 2 {
		channelPtr := *(**C.int16_t)(unsafe.Add(unsafe.Pointer(b), uintptr(i)*unsafe.Sizeof(*b)))
		for j := range numSamples {
			d.samples[i] = append(d.samples[i], int16(*(*C.int16_t)(unsafe.Add(unsafe.Pointer(channelPtr), uintptr(j)*C.sizeof_int16_t))))
		}
		C.av_freep(unsafe.Pointer(&channelPtr))
	}
	C.av_freep(unsafe.Pointer(&b))
}

func (d *AVDecoder) Position() int {
	return d.pos / avBytesPerFrame
}

func (d *AVDecoder) Len() int {
	return int(d.NbSamples())
}

func (d *AVDecoder) Close() error {
	d.free()
	d.samples = [2][]int16{}
	d.pos = 0
	return nil
}

func (d *AVDecoder) Err() error {
	return d.err
}

func (d *AVDecoder) Seek(p int) error {
	if p < 0 || d.Len() < p {
		return fmt.Errorf("seek position %v out of range: [%v, %v]", p, 0, d.Len())
	}
	ret := C.decoder_seek(d.dec, C.int64_t(p))
	if ret < 0 {
		return ErrSeek
	}
	d.pos = p * avBytesPerFrame
	return nil
}
