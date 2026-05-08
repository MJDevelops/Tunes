package audio

// #include <stdlib.h>
// #include <stdint.h>
// #include "decode.h"
// #include "libavutil/avutil.h"
import "C"
import (
	"unsafe"
)

type AVDecoder struct {
	buf         **C.double_t
	bufArr      [2][]float64
	channelSize int64
}

// TODO: simplify?
func (d *AVDecoder) DecodeAudio(filename string) [2][]float64 {
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(file))
	d.buf = (**C.double_t)(C.av_malloc(2 * (C.size_t)(unsafe.Sizeof(*(*d.buf)))))
	defer d.freeBuffer()
	d.channelSize = int64(C.decode(d.buf, file))

	for i := range 2 {
		channelPtr := *(**C.double_t)(unsafe.Add(unsafe.Pointer(d.buf), uintptr(i)*unsafe.Sizeof(*d.buf)))
		d.bufArr[i] = make([]float64, d.channelSize)
		for j := range d.channelSize {
			d.bufArr[i][j] = float64(*(*C.double_t)(unsafe.Add(unsafe.Pointer(channelPtr), uintptr(j)*C.sizeof_double_t)))
		}
	}
	return d.bufArr
}

func (d *AVDecoder) freeBuffer() {
	C.free_sample_buffer((*unsafe.Pointer)(unsafe.Pointer(d.buf)), 2, C.int64_t(d.channelSize))
}
