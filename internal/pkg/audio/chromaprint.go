package audio

// #include <stdlib.h>
// #include "fingerprint.h"
import "C"
import (
	"unsafe"

	"github.com/mjdevelops/tunes/internal/pkg/os"
)

func FingerprintFile(path string) (string, error) {
	if isFile := os.IsFile(path); !isFile {
		return "", ErrNotAFile
	}

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	fingerprint := C.fingerprint_file(cpath)
	defer C.free(unsafe.Pointer(fingerprint))

	return C.GoString(fingerprint), nil
}
