package audio

// #include <stdlib.h>
// #include <stdint.h>
// #include "decode.h"
// #include "libavutil/avutil.h"
import "C"
import (
	"unsafe"
)

type Decoded struct {
	buf *C.int16_t
}

func (d *Decoded) DecodeAudio(filename string) []int16 {
	file := C.CString(filename)
	defer C.free(unsafe.Pointer(file))
	d.buf = C.decode(file)
	return ((*[1 << 30]int16)(unsafe.Pointer(d.buf)))[:]
}

func (d *Decoded) FreeBuffer() {
	C.av_freep(unsafe.Pointer(d.buf))
}
