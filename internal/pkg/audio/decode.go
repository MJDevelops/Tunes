package audio

// #include <stdlib.h>
// #include <stdint.h>
// #include "decode.h"
// #include "libavutil/avutil.h"
import "C"
import (
	"encoding/binary"
	"errors"
	"unsafe"

	"github.com/gopxl/beep/v2"
)

const (
	avPrecision     = 2
	avChannels      = 2
	avBytesPerFrame = avChannels * avPrecision
)

type AVDecoder struct {
	buf     *C.SampleBuffer
	f       beep.Format
	pos     int
	file    string
	err     error
	samples [2][]int16
}

var ErrDecoding = errors.New("error during decoding")

// TODO: Implement proper streaming of audio files and
// don't decode the file in one go as its slow
func NewDecoder(filename string) (s beep.StreamSeekCloser, format beep.Format, err error) {
	d := &AVDecoder{}
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(file))
	d.buf = C.sb_alloc()
	if d.buf == nil {
		return nil, beep.Format{}, ErrDecoding
	}

	ret := C.decode(d.buf, file)
	if ret < 0 {
		return nil, beep.Format{}, ErrDecoding
	}

	for i := range 2 {
		channelPtr := *(**C.int16_t)(unsafe.Add(unsafe.Pointer(d.buf.data), uintptr(i)*unsafe.Sizeof(*d.buf.data)))
		d.samples[i] = make([]int16, d.buf.channel_size)
		for j := range d.buf.channel_size {
			d.samples[i][j] = int16(*(*C.int16_t)(unsafe.Add(unsafe.Pointer(channelPtr), uintptr(j)*C.sizeof_int16_t)))
		}
	}

	format = beep.Format{
		Precision:   avPrecision,
		NumChannels: avChannels,
		SampleRate:  beep.SampleRate(d.SampleRate()),
	}

	d.f = format

	return d, format, nil
}

func (d *AVDecoder) Samples() *[2][]int16 {
	return &d.samples
}

func (d *AVDecoder) SampleRate() int {
	if d.buf != nil {
		return int(d.buf.sample_rate)
	}
	return 0
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
}

func (d *AVDecoder) Stream(samples [][2]float64) (n int, ok bool) {
	b := make([]byte, avBytesPerFrame)
	for i := range samples {
		pos := d.Position()
		if pos >= len(d.samples[0]) {
			break
		}

		// ffmpeg uses native endianness for audio samples
		binary.NativeEndian.PutUint16(b, uint16(d.samples[0][pos]))
		binary.NativeEndian.PutUint16(b[2:], uint16(d.samples[1][pos]))

		samples[i], _ = d.f.DecodeSigned(b)
		d.pos += avBytesPerFrame
		n++
		ok = true
	}
	return n, ok
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
	// TODO
	return nil
}
