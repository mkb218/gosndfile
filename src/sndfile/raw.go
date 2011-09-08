package sndfile

// #cgo LDFLAGS: -lsndfile
// #import <sndfile.h>
import "C"

import (
	"unsafe"
	"os"
)

// needs test, needs doc
func (f *File) ReadRaw(data []byte) (read int64, err os.Error) {
	read = int64(C.sf_read_raw(f.s, unsafe.Pointer(&data[0]), C.sf_count_t(len(data))))
	if read != int64(len(data)) {
		err = sErrorType(C.sf_error(f.s))
	}
	return
}

// needs test, needs doc
func (f *File) WriteRaw(data []byte) (written int64, err os.Error) {
	written = int64(C.sf_write_raw(f.s, unsafe.Pointer(&data[0]), C.sf_count_t(len(data))))
	if written != int64(len(data)) {
		err = sErrorType(C.sf_error(f.s))
	}
	return
}
