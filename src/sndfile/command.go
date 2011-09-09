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
	return f.genericBoolBoolCmd(C.SFC_SET_NORM_FLOAT, norm)
}

/*This command only affects data read from or written to using ReadItems, ReadFrames, WriteItems, or WriteFrames with slices of float64.

For read operations setting normalisation to true means that the data from all subsequent reads will be be normalised to the range [-1.0, 1.0].

For write operations, setting normalisation to true means than all data supplied to the float write functions should be in the range [-1.0, 1.0] and will be scaled for the file format as necessary.

For both cases, setting normalisation to false means that no scaling will take place.

Returns the previous normalization setting. */
//needs test
func (f *File) SetDoubleNormalization(norm bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_NORM_DOUBLE, norm)
}

// Returns the current float32 normalization mode.
//needs test
func (f *File) GetFloatNormalization() bool {
	return f.genericBoolBoolCmd(C.SFC_GET_NORM_FLOAT, false)
}

// Returns the current float64 normalization mode.
//needs test
func (f *File) GetDoubleNormalization() bool {
	return f.genericBoolBoolCmd(C.SFC_GET_NORM_DOUBLE, false)
}

//Set/clear the scale factor when integer (short/int) data is read from a file containing floating point data.
//needs test
func (f *File) SetFloatIntScaleRead(scale bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_SCALE_FLOAT_INT_READ, scale)
}

//Set/clear the scale factor when integer (short/int) data is written to a file as floating point data.
//needs test
func (f *File) SetIntFloatScaleWrite(scale bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_SCALE_INT_FLOAT_WRITE, scale)
}

//Retrieve the number of simple formats supported by libsndfile.
//needstest
func GetSimpleFormatCount() int {
	var o C.int
	C.sf_command(nil, C.SFC_GET_SIMPLE_FORMAT_COUNT, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o)))
	return int(o)
}

//Retrieve information about a simple format.
//The value of the format argument should be the format number (ie 0 <= format <= count value obtained using GetSimpleFormatCount()).
// The returned format argument is suitable for use in sndfile.Open()
//needs test , needs example, needs doc
func GetSimpleFormat(format int) (oformat int, name string, extension string, ok bool) {
	var o C.SF_FORMAT_INFO
	o.format = C.int(format)
	ok = (0 == C.sf_command(nil, C.SFC_GET_SIMPLE_FORMAT, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o))))
	oformat = int(o.format)
	name = C.GoString(o.name)
	extension = C.GoString(o.extension)
	return
}

//When GetFormatInfo() is called, the format argument is examined and if (format & SF_FORMAT_TYPEMASK) is a valid format then the returned strings contain information about the given major type. If (format & SF_FORMAT_TYPEMASK) is FALSE and (format & SF_FORMAT_SUBMASK) is a valid subtype format then the returned strings contain information about the given subtype.
//needs test , needs example, needs doc
func GetFormatInfo(format int) (oformat int, name string, extension string, ok bool) {
	var o C.SF_FORMAT_INFO
	o.format = C.int(format)
	ok = (0 == C.sf_command(nil, C.SFC_GET_FORMAT_INFO, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o))))
	oformat = int(o.format)
	name = C.GoString(o.name)
	extension = C.GoString(o.extension)
	return
}

//Retrieve the number of major formats supported by libsndfile.
//needstest
func GetMajorFormatCount() int {
	var o C.int
	C.sf_command(nil, C.SFC_GET_FORMAT_MAJOR_COUNT, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o)))
	return int(o)
}

//Retrieve information about a major format type
//For a more comprehensive example, see the program list_formats.c in the examples/ directory of the libsndfile source code distribution.
//needs test , needs example, needs doc
func GetMajorFormatInfo(format int) (oformat int, name string, extension string, ok bool) {
	var o C.SF_FORMAT_INFO
	o.format = C.int(format)
	ok = (0 == C.sf_command(nil, C.SFC_GET_FORMAT_MAJOR, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o))))
	oformat = int(o.format)
	name = C.GoString(o.name)
	extension = C.GoString(o.extension)
	return
}

//Retrieve the number of subformats supported by libsndfile.
//needstest
func GetSubFormatCount() int {
	var o C.int
	C.sf_command(nil, C.SFC_GET_FORMAT_SUBTYPE_COUNT, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o)))
	return int(o)
}

//Enumerate the subtypes (this function does not translate a subtype into a string describing that subtype). A typical use case might be retrieving a string description of all subtypes so that a dialog box can be filled in.
//needs test , needs example, needs doc
func GetSubFormatInfo(format int) (oformat int, name string, ok bool) {
	var o C.SF_FORMAT_INFO
	o.format = C.int(format)
	ok = (0 == C.sf_command(nil, C.SFC_GET_FORMAT_SUBTYPE, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o))))
	oformat = int(o.format)
	name = C.GoString(o.name)
	return
}

//By default, WAV and AIFF files which contain floating point data (subtype SF_FORMAT_FLOAT or SF_FORMAT_DOUBLE) have a PEAK chunk. By using this command, the addition of a PEAK chunk can be turned on or off.

//Note : This call must be made before any data is written to the file.
// needstest
func (f *File) SetAddPeakChunk(set bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_ADD_PEAK_CHUNK, set)
}

//The header of an audio file is normally written by libsndfile when the file is closed using sf_close().

//There are however situations where large files are being generated and it would be nice to have valid data in the header before the file is complete. Using this command will update the file header to reflect the amount of data written to the file so far. Other programs opening the file for read (before any more data is written) will then read a valid sound file header.
//needs test
func (f *File) UpdateHeaderNow() {
	C.sf_command(f.s, C.SFC_UPDATE_HEADER_NOW, nil, 0)
}

//Similar to SFC_UPDATE_HEADER_NOW but updates the header at the end of every call to the sf_write* functions.
//needstest
func (f *File) SetUpdateHeaderAuto(set bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_UPDATE_HEADER_AUTO, set)
}

func (f *File) genericBoolBoolCmd(cmd C.int, i bool) bool {
	ib := C.SF_FALSE
	if i { 
		ib = C.SF_TRUE
	}
	
	n := C.sf_command(f.s, cmd, nil, C.int(ib))
	return (n == C.SF_TRUE)
}

// This allows libsndfile experts to use the command interface for commands not currently supported. See http://www.mega-nerd.com/libsndfile/command.html
// The f argument may be nil in cases where the command does not require a SNDFILE argument.
// The method's cmd, data, and datasize arguments are used the same way as the correspondingly named arguments for sf_command
//needs test
func GenericCmd(f *File, cmd C.int, data unsafe.Pointer, datasize int) int {
	var s *C.SNDFILE = nil
	if f != nil {
		s = f.s
	}
	return int(C.sf_command(s, cmd, data, C.int(datasize)))
}