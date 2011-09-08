package sndfile

// #cgo LDFLAGS: -lsndfile
// #include <stdlib.h>
// #include <sndfile.h>
import "C"

import "unsafe"
import "os"
import "fmt"

// GetLibVersion retrieves the version of the library as a string
func GetLibVersion() (s string, err os.Error) {
	l := C.sf_command(nil, C.SFC_GET_LIB_VERSION, nil, 0)
	c := make([]byte, l)
	m := C.sf_command(nil, C.SFC_GET_LIB_VERSION, unsafe.Pointer(&c[0]), l)

	if m != l {
		err = os.NewError(fmt.Sprintf("GetLibVersion: expected %d bytes in string, recv'd %d", l, m))
	}
	s = string(c)
	return
}

// Retrieve the log buffer generated when opening a file as a string. This log buffer can often contain a good reason for why libsndfile failed to open a particular file.
//needs test
func (f *File) GetLogInfo() (s string, err os.Error) {
	l := C.sf_command(f.s, C.SFC_GET_LOG_INFO, nil, 0)
	c := make([]byte, l)
	m := C.sf_command(f.s, C.SFC_GET_LOG_INFO, unsafe.Pointer(&c[0]), l)

	if m != l {
		err = os.NewError(fmt.Sprintf("GetLogInfo: expected %d bytes in string, recv'd %d", l, m))
	}
	s = string(c)
	return
}

// Retrieve the measured maximum signal value. This involves reading through the whole file which can be slow on large files.
//needs test
func (f *File) CalcSignalMax() (ret float64, err os.Error) {
	e := C.sf_command(f.s, C.SFC_CALC_SIGNAL_MAX, unsafe.Pointer(&ret), 8)
	if e != 0 {
		err = sErrorType(e)
	}
	return
}

// Retrieve the measured normalised maximum signal value. This involves reading through the whole file which can be slow on large files.
//needs test
func (f *File) CalcNormSignalMax() (ret float64, err os.Error) {
	e := C.sf_command(f.s, C.SFC_CALC_NORM_SIGNAL_MAX, unsafe.Pointer(&ret), 8)
	if e != 0 {
		err = sErrorType(e)
	}
	return
}

//Calculate the peak value (ie a single number) for each channel. This involves reading through the whole file which can be slow on large files.
//needs test
func (f *File) CalcMaxAllChannels() (ret []float64, err os.Error) {
	c := f.Format.Channels
	ret = make([]float64, c)
	e := C.sf_command(f.s, C.SFC_CALC_MAX_ALL_CHANNELS, unsafe.Pointer(&ret[0]), C.int(c*8))
	if e != 0 {
		err = sErrorType(e)
	}
	return
}

//Calculate the normalised peak for each channel. This involves reading through the whole file which can be slow on large files.
//needs test
func (f *File) CalcNormMaxAllChannels() (ret []float64, err os.Error) {
	c := f.Format.Channels
	ret = make([]float64, c)
	e := C.sf_command(f.s, C.SFC_CALC_NORM_MAX_ALL_CHANNELS, unsafe.Pointer(&ret[0]), C.int(c*8))
	if e != 0 {
		err = sErrorType(e)
	}
	return
}

//Retrieve the peak value for the file as stored in the file header.
//needs test
func (f *File) GetSignalMax() (ret float64, ok bool) {
	r := C.sf_command(f.s, C.SFC_GET_SIGNAL_MAX, unsafe.Pointer(&ret), 8)
	if r == C.SF_TRUE {
		ok = true
	}
	return
}

//Retrieve the peak value for the file as stored in the file header.
//needs test
func (f *File) GetMaxAllChannels() (ret []float64, ok bool) {
	c := f.Format.Channels
	ret = make([]float64, c)
	e := C.sf_command(f.s, C.SFC_GET_MAX_ALL_CHANNELS, unsafe.Pointer(&ret[0]), C.int(c*8))
	if e == C.SF_TRUE {
		ok = true
	}
	return
}

/*This command only affects data read from or written to using ReadItems, ReadFrames, WriteItems, or WriteFrames with slices of float32.

For read operations setting normalisation to true means that the data from all subsequent reads will be be normalised to the range [-1.0, 1.0].

For write operations, setting normalisation to true means than all data supplied to the float write functions should be in the range [-1.0, 1.0] and will be scaled for the file format as necessary.

For both cases, setting normalisation to false means that no scaling will take place.

Returns the previous normalization setting. */
//needs test
func (f *File) SetFloatNormalization(norm bool) bool {
	i := C.SF_FALSE
	if norm {
		i = C.SF_TRUE
	}
	n := C.sf_command(f.s, C.SFC_SET_NORM_FLOAT, nil, C.int(i))
	if n == C.SF_TRUE {
		return true
	}
	return false
}

/*This command only affects data read from or written to using ReadItems, ReadFrames, WriteItems, or WriteFrames with slices of float64.

For read operations setting normalisation to true means that the data from all subsequent reads will be be normalised to the range [-1.0, 1.0].

For write operations, setting normalisation to true means than all data supplied to the float write functions should be in the range [-1.0, 1.0] and will be scaled for the file format as necessary.

For both cases, setting normalisation to false means that no scaling will take place.

Returns the previous normalization setting. */
//needs test
func (f *File) SetDoubleNormalization(norm bool) bool {
	i := C.SF_FALSE
	if norm {
		i = C.SF_TRUE
	}
	n := C.sf_command(f.s, C.SFC_SET_NORM_DOUBLE, nil, C.int(i))
	if n == C.SF_TRUE {
		return true
	}
	return false
}

// Returns the current float32 normalization mode.
//needs test
func (f *File) GetFloatNormalization() bool {
	return (C.sf_command(f.s, C.SFC_GET_NORM_FLOAT, nil, 0) == C.SF_TRUE)
}

// Returns the current float64 normalization mode.
//needs test
func (f *File) GetDoubleNormalization() bool {
	return (C.sf_command(f.s, C.SFC_GET_NORM_DOUBLE, nil, 0) == C.SF_TRUE)
}

// oh god this is boring.
//needs test, needs doc
func (f *File) GenericCmd(cmd C.int, data unsafe.Pointer, datasize int) int {
	return int(C.sf_command(f.s, cmd, data, C.int(datasize)))
}