package audio

// #include <stdlib.h>
// #include <stdint.h>
// #include "decode.h"
// #include "libavutil/avutil.h"
import "C"
import (
	"errors"
	"unsafe"
)

type AVDecoder struct {
	buf     *C.SampleBuffer
	samples [2][]float64
}

var ErrDecoding = errors.New("error during decoding")

func (d *AVDecoder) DecodeAudio(filename string) error {
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(file))
	d.buf = C.sb_alloc()
	ret := int(C.decode(d.buf, file))

	if ret < 0 {
		return ErrDecoding
	}

	for i := range 2 {
		channelPtr := *(**C.double_t)(unsafe.Add(unsafe.Pointer(d.buf.data), uintptr(i)*unsafe.Sizeof(*d.buf.data)))
		d.samples[i] = make([]float64, d.buf.channel_size)
		for j := range d.buf.channel_size {
			d.samples[i][j] = float64(*(*C.double_t)(unsafe.Add(unsafe.Pointer(channelPtr), uintptr(j)*C.sizeof_double_t)))
		}
	}

	return nil
}

func (d *AVDecoder) Samples() [2][]float64 {
	return d.samples
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

func (d *AVDecoder) Free() {
	if d.buf != nil {
		C.sb_free(d.buf)
	}
}
