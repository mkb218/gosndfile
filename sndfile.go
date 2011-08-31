package sndfile

// The sndfile package is a binding for libsndfile. It packages the libsndfile API in a go-like manner.

// #cgo CFLAGS: -I/opt/local/include
// #cgo LDFLAGS: -L/opt/local/lib -lsndfile
// #include <stdlib.h>
// #include <sndfile.h>
import "C"

import (
	"os"
	"unsafe"
	"reflect"
)

// A sound file. Does not conform to io.Reader.
type File struct {
	s *C.SNDFILE
	i Info
}

// sErrorType represents a sndfile API error and grabs error description strings from the API.
type sErrorType C.int

func (e sErrorType) String() string {
	return C.GoString(C.sf_error_number(C.int(e)))
}

// File mode: Read, Write, or ReadWrite
type Mode int

const (
	Read      Mode = C.SFM_READ
	Write     Mode = C.SFM_WRITE
	ReadWrite Mode = C.SFM_RDWR
)

// Info is the struct needed to open a file for reading or writing. When opening a file for reading, everything may generally be left zeroed. The only exception to this is the case of RAW files where the caller has to set the samplerate, channels and format fields to valid values.
type Info struct {
	Frames       int64
	Samplerate   int32
	Channels     int32
	Format       int32
	Sections     int32
	Seekable     int32
	Pad_godefs_0 [4]byte
}

// The format field in the above SF_INFO structure is made up of the bit-wise OR of a major format type (values between 0x10000 and 0x08000000), a minor format type (with values less than 0x10000) and an optional endian-ness value. The currently understood formats are listed in sndfile.h as follows and also include bitmasks for separating major and minor file types. Not all combinations of endian-ness and major and minor file types are valid.
type Format int32

const (
	SF_FORMAT_WAV   Format = 0x010000 /* Microsoft WAV format (little endian). */
	SF_FORMAT_AIFF  Format = 0x020000 /* Apple/SGI AIFF format (big endian). */
	SF_FORMAT_AU    Format = 0x030000 /* Sun/NeXT AU format (big endian). */
	SF_FORMAT_RAW   Format = 0x040000 /* RAW PCM data. */
	SF_FORMAT_PAF   Format = 0x050000 /* Ensoniq PARIS file format. */
	SF_FORMAT_SVX   Format = 0x060000 /* Amiga IFF / SVX8 / SV16 format. */
	SF_FORMAT_NIST  Format = 0x070000 /* Sphere NIST format. */
	SF_FORMAT_VOC   Format = 0x080000 /* VOC files. */
	SF_FORMAT_IRCAM Format = 0x0A0000 /* Berkeley/IRCAM/CARL */
	SF_FORMAT_W64   Format = 0x0B0000 /* Sonic Foundry's 64 bit RIFF/WAV */
	SF_FORMAT_MAT4  Format = 0x0C0000 /* Matlab (tm) V4.2 / GNU Octave 2.0 */
	SF_FORMAT_MAT5  Format = 0x0D0000 /* Matlab (tm) V5.0 / GNU Octave 2.1 */
	SF_FORMAT_PVF   Format = 0x0E0000 /* Portable Voice Format */
	SF_FORMAT_XI    Format = 0x0F0000 /* Fasttracker 2 Extended Instrument */
	SF_FORMAT_HTK   Format = 0x100000 /* HMM Tool Kit format */
	SF_FORMAT_SDS   Format = 0x110000 /* Midi Sample Dump Standard */
	SF_FORMAT_AVR   Format = 0x120000 /* Audio Visual Research */
	SF_FORMAT_WAVEX Format = 0x130000 /* MS WAVE with WAVEFORMATEX */
	SF_FORMAT_SD2   Format = 0x160000 /* Sound Designer 2 */
	SF_FORMAT_FLAC  Format = 0x170000 /* FLAC lossless file format */
	SF_FORMAT_CAF   Format = 0x180000 /* Core Audio File format */
	SF_FORMAT_WVE   Format = 0x190000 /* Psion WVE format */
	SF_FORMAT_OGG   Format = 0x200000 /* Xiph OGG container */
	SF_FORMAT_MPC2K Format = 0x210000 /* Akai MPC 2000 sampler */
	SF_FORMAT_RF64  Format = 0x220000 /* RF64 WAV file */

	/* Subtypes from here on. */

	SF_FORMAT_PCM_S8 Format = 0x0001 /* Signed 8 bit data */
	SF_FORMAT_PCM_16 Format = 0x0002 /* Signed 16 bit data */
	SF_FORMAT_PCM_24 Format = 0x0003 /* Signed 24 bit data */
	SF_FORMAT_PCM_32 Format = 0x0004 /* Signed 32 bit data */

	SF_FORMAT_PCM_U8 Format = 0x0005 /* Unsigned 8 bit data (WAV and RAW only) */

	SF_FORMAT_FLOAT  Format = 0x0006 /* 32 bit float data */
	SF_FORMAT_DOUBLE Format = 0x0007 /* 64 bit float data */

	SF_FORMAT_ULAW      Format = 0x0010 /* U-Law encoded. */
	SF_FORMAT_ALAW      Format = 0x0011 /* A-Law encoded. */
	SF_FORMAT_IMA_ADPCM Format = 0x0012 /* IMA ADPCM. */
	SF_FORMAT_MS_ADPCM  Format = 0x0013 /* Microsoft ADPCM. */

	SF_FORMAT_GSM610    Format = 0x0020 /* GSM 6.10 encoding. */
	SF_FORMAT_VOX_ADPCM Format = 0x0021 /* Oki Dialogic ADPCM encoding. */

	SF_FORMAT_G721_32 Format = 0x0030 /* 32kbs G721 ADPCM encoding. */
	SF_FORMAT_G723_24 Format = 0x0031 /* 24kbs G723 ADPCM encoding. */
	SF_FORMAT_G723_40 Format = 0x0032 /* 40kbs G723 ADPCM encoding. */

	SF_FORMAT_DWVW_12 Format = 0x0040 /* 12 bit Delta Width Variable Word encoding. */
	SF_FORMAT_DWVW_16 Format = 0x0041 /* 16 bit Delta Width Variable Word encoding. */
	SF_FORMAT_DWVW_24 Format = 0x0042 /* 24 bit Delta Width Variable Word encoding. */
	SF_FORMAT_DWVW_N  Format = 0x0043 /* N bit Delta Width Variable Word encoding. */

	SF_FORMAT_DPCM_8  Format = 0x0050 /* 8 bit differential PCM (XI only) */
	SF_FORMAT_DPCM_16 Format = 0x0051 /* 16 bit differential PCM (XI only) */

	SF_FORMAT_VORBIS Format = 0x0060 /* Xiph Vorbis encoding. */

	/* Endian-ness options. */

	SF_ENDIAN_FILE   Format = 0x00000000 /* Default file endian-ness. */
	SF_ENDIAN_LITTLE Format = 0x10000000 /* Force little endian-ness. */
	SF_ENDIAN_BIG    Format = 0x20000000 /* Force big endian-ness. */
	SF_ENDIAN_CPU    Format = 0x30000000 /* Force CPU endian-ness. */

	SF_FORMAT_SUBMASK  Format = 0x0000FFFF
	SF_FORMAT_TYPEMASK Format = 0x0FFF0000
	SF_FORMAT_ENDMASK  Format = 0x30000000
)

func Open(name string, mode Mode, info Info) (o File, err os.Error) {
	c := C.CString(name)
	defer C.free(unsafe.Pointer(c))
	o.s = C.sf_open(c, C.int(mode), (*C.SF_INFO)(unsafe.Pointer(&info)))
	if o.s == nil {
		err = sErrorType(C.sf_error(o.s))
	}
	o.i = info
	return
}

// This probably won't work on windows
func OpenFd(fd int, mode Mode, info Info, close_desc bool) (o File, err os.Error) {
	var cd C.int
	if close_desc {
		cd = 1
	}
	o.s = C.sf_open_fd(C.int(fd), C.int(mode), (*C.SF_INFO)(unsafe.Pointer(&info)), cd)
	if o.s == nil {
		err = sErrorType(C.sf_error(o.s))
	}
	o.i = info
	return
}

// not interested in dealing with callbacks from c to go right now kthx, so no sf_open_virtual

// This function allows the caller to check if a set of parameters in the Info struct is valid before calling Open in Write mode.
// FormatCheck returns true if the parameters are valid and false otherwise. */
func FormatCheck(i Info) bool {
	return C.sf_format_check((*C.SF_INFO)(unsafe.Pointer(&i))) != C.SF_TRUE
}

// Whence args for Seek()
type Whence C.int

const (
	Set     Whence = C.SEEK_SET // The offset is set to the start of the audio data plus offset (multichannel) frames.
	Current Whence = C.SEEK_CUR //The offset is set to its current location plus offset (multichannel) frames.
	End     Whence = C.SEEK_END //The offset is set to the end of the data plus offset (multichannel) frames.
)

//The file seek functions work much like lseek in unistd.h with the exception that the non-audio data is ignored and the seek only moves within the audio data section of the file. In addition, seeks are defined in number of (multichannel) frames. Therefore, a seek in a stereo file from the current position forward with an offset of 1 would skip forward by one sample of both channels. This function returns the new offset, and a non-nil error value if unsuccessful
func (f File) Seek(frames int64, w Whence) (offset int64, err os.Error) {
	r := C.sf_seek(f.s, C.sf_count_t(frames), C.int(w))
	if r == -1 {
		err = sErrorType(C.sf_error(f.s))
	} else {
		offset = int64(r)
	}
	return
}

// The close function closes the file, deallocates its internal buffers and returns a non-nil error value in cas of error
func (f File) Close() (err os.Error) {
	if C.sf_close(f.s) != 0 {
		err = sErrorType(C.sf_error(f.s))
	}
	return
}

//If the file is opened Write or ReadWrite, call the operating system's function to force the writing of all file cache buffers to disk. If the file is opened Read no action is taken.
func (f File) WriteSync() {
	C.sf_write_sync(f.s)
}

/*The file read items functions fill the array pointed to by out with the requested number of items. The length of out must be a product of the number of channels or an error will occur.

It is important to note that the data type used by the calling program and the data format of the file do not need to be the same. For instance, it is possible to open a 16 bit PCM encoded WAV file and read the data into a slice of floats. The library seamlessly converts between the two formats on-the-fly. See Note 1.

Returns the number of items read. Unless the end of the file was reached during the read, the return value should equal the number of items requested. Attempts to read beyond the end of the file will not result in an error but will cause ReadItems to return less than the number of items requested or 0 if already at the end of the file.*/
func (f File) ReadItems(out []interface{}) (read int64, err os.Error) {
	items := len(out)
	if items < 1 {
		err = os.EOF
		return
	}
	var n C.sf_count_t
	t := reflect.TypeOf(out[0])
	switch t.Kind() {
	case reflect.Int16:
		fallthrough
	case reflect.Uint16:
		n = C.sf_read_short(f.s, (*C.short)(unsafe.Pointer(&out[0])), C.sf_count_t(items))
	case reflect.Int32:
		fallthrough
	case reflect.Uint32:
		n = C.sf_read_int(f.s, (*C.int)(unsafe.Pointer(&out[0])), C.sf_count_t(items))
	case reflect.Float32:
		n = C.sf_read_float(f.s, (*C.float)(unsafe.Pointer(&out[0])), C.sf_count_t(items))
	case reflect.Float64:
		n = C.sf_read_double(f.s, (*C.double)(unsafe.Pointer(&out[0])), C.sf_count_t(items))
	case reflect.Int:
		fallthrough
	case reflect.Uint:
		switch t.Bits() {
		case 32:
			n = C.sf_read_int(f.s, (*C.int)(unsafe.Pointer(&out[0])), C.sf_count_t(items))
		case 16: // doubtful
			n = C.sf_read_short(f.s, (*C.short)(unsafe.Pointer(&out[0])), C.sf_count_t(items))
		default:
			err = os.NewError("Unsupported type in read buffer, needs (u)int16, (u)int32, or float type")
		}
	default:
		err = os.NewError("Unsupported type in read buffer, needs (u)int16, (u)int32, or float type")
	}
	if err != nil {
		read = -1
		return
	}

	read = int64(n)
	if read < 0 {
		err = sErrorType(C.sf_error(f.s))
	}
	return
}

/*The file read frames functions fill the array pointed to by out with the requested number of frames of data. The array must be large enough to hold the product of frames and the number of channels.

The sf_readf_XXXX functions return the number of frames read. Unless the end of the file was reached during the read, the return value should equal the number of frames requested. Attempts to read beyond the end of the file will not result in an error but will cause the sf_readf_XXXX functions to return less than the number of frames requested or 0 if already at the end of the file.*/
func (f File) ReadFrames(out []interface{}) (read int64, err os.Error) {
	frames := len(out)/int(f.i.Channels)
	if frames < 1 {
		err = os.EOF
		return
	}
	var n C.sf_count_t
	t := reflect.TypeOf(out[0])
	switch t.Kind() {
	case reflect.Int16:
		fallthrough
	case reflect.Uint16:
		n = C.sf_readf_short(f.s, (*C.short)(unsafe.Pointer(&out[0])), C.sf_count_t(frames))
	case reflect.Int32:
		fallthrough
	case reflect.Uint32:
		n = C.sf_readf_int(f.s, (*C.int)(unsafe.Pointer(&out[0])), C.sf_count_t(frames))
	case reflect.Float32:
		n = C.sf_readf_float(f.s, (*C.float)(unsafe.Pointer(&out[0])), C.sf_count_t(frames))
	case reflect.Float64:
		n = C.sf_readf_double(f.s, (*C.double)(unsafe.Pointer(&out[0])), C.sf_count_t(frames))
	case reflect.Int:
		fallthrough
	case reflect.Uint:
		switch t.Bits() {
		case 32:
			n = C.sf_readf_int(f.s, (*C.int)(unsafe.Pointer(&out[0])), C.sf_count_t(frames))
		case 16: // doubtful
			n = C.sf_readf_short(f.s, (*C.short)(unsafe.Pointer(&out[0])), C.sf_count_t(frames))
		default:
			err = os.NewError("Unsupported type in read buffer, needs (u)int16, (u)int32, or float type")
		}
	default:
		err = os.NewError("Unsupported type in read buffer, needs (u)int16, (u)int32, or float type")
	}
	if err != nil {
		read = -1
		return
	}

	read = int64(n)
	if read < 0 {
		err = sErrorType(C.sf_error(f.s))
	}
	return
}

// going to skip raw I/O for now