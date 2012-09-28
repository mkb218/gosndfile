package sndfile

// #cgo pkg-config: sndfile
// #include <sndfile.h>
import "C"

import "unsafe"

//Note: Unless you are writing an external decoder/encode that uses libsndfile to handle the file headers, you should not be using this function.

//The raw read and write functions read raw audio data from the audio file (not to be confused with reading RAW header-less PCM files). The number of bytes read or written must always be an integer multiple of the number of channels multiplied by the number of bytes required to represent one sample from one channel.

//The raw read and write functions return the number of bytes read or written (which should be the same as the bytes parameter) and any error that occurs while reading or writing
// needs test
func (f *File) ReadRaw(data []byte) (read int64, err error) {
	read = int64(C.sf_read_raw(f.s, unsafe.Pointer(&data[0]), C.sf_count_t(len(data))))
	if read != int64(len(data)) {
		err = sErrorType(C.sf_error(f.s))
	}
	return
}

//Note: Unless you are writing an external decoder/encode that uses libsndfile to handle the file headers, you should not be using this function.

//The raw read and write functions read raw audio data from the audio file (not to be confused with reading RAW header-less PCM files). The number of bytes read or written must always be an integer multiple of the number of channels multiplied by the number of bytes required to represent one sample from one channel.

//The raw read and write functions return the number of bytes read or written (which should be the same as the bytes parameter) and any error that occurs while reading or writing
// needs test
func (f *File) WriteRaw(data []byte) (written int64, err error) {
	written = int64(C.sf_write_raw(f.s, unsafe.Pointer(&data[0]), C.sf_count_t(len(data))))
	if written != int64(len(data)) {
		err = sErrorType(C.sf_error(f.s))
	}
	return
}
