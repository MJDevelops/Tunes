package audio

// #include <stdlib.h>
// #include "fingerprint.h"
import "C"
import "unsafe"

func FingerprintFile(path string) string {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	fingerprint := C.fingerprint_file(cpath)
	defer C.free(unsafe.Pointer(fingerprint))

	return C.GoString(fingerprint)
}
