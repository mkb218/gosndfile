package sndfile

// #cgo pkg-config: sndfile
// #include <stdlib.h>
// #include <sndfile.h>
// #include <string.h>
import "C"

import (
	"errors"
	"strings"
)
import "unsafe"

import "fmt"

// GetLibVersion retrieves the version of the library as a string
func GetLibVersion() (s string, err error) {
	l := C.sf_command(nil, C.SFC_GET_LIB_VERSION, nil, 0)
	c := make([]byte, l)
	m := C.sf_command(nil, C.SFC_GET_LIB_VERSION, unsafe.Pointer(&c[0]), l)

	if m != l {
		err = errors.New(fmt.Sprintf("GetLibVersion: expected %d bytes in string, recv'd %d", l, m))
	}
	s = string(c)
	return
}

// Retrieve the log buffer generated when opening a file as a string. This log buffer can often contain a good reason for why libsndfile failed to open a particular file.
func (f *File) GetLogInfo() (s string, err error) {
	l := C.sf_command(f.s, C.SFC_GET_LOG_INFO, nil, 0)
	c := make([]byte, l)
	m := C.sf_command(f.s, C.SFC_GET_LOG_INFO, unsafe.Pointer(&c[0]), l)

	if m != l {
		c = c[0:m]
	}
	s = string(c)
	return
}

// Retrieve the measured maximum signal value. This involves reading through the whole file which can be slow on large files.
func (f *File) CalcSignalMax() (ret float64, err error) {
	e := C.sf_command(f.s, C.SFC_CALC_SIGNAL_MAX, unsafe.Pointer(&ret), 8)
	if e != 0 {
		err = sErrorType(e)
	}
	return
}

// Retrieve the measured normalised maximum signal value. This involves reading through the whole file which can be slow on large files.
func (f *File) CalcNormSignalMax() (ret float64, err error) {
	e := C.sf_command(f.s, C.SFC_CALC_NORM_SIGNAL_MAX, unsafe.Pointer(&ret), 8)
	if e != 0 {
		err = sErrorType(e)
	}
	return
}

//Calculate the peak value (ie a single number) for each channel. This involves reading through the whole file which can be slow on large files.
func (f *File) CalcMaxAllChannels() (ret []float64, err error) {
	c := f.Format.Channels
	ret = make([]float64, c)
	e := C.sf_command(f.s, C.SFC_CALC_MAX_ALL_CHANNELS, unsafe.Pointer(&ret[0]), C.int(c*8))
	if e != 0 {
		err = sErrorType(e)
	}
	return
}

//Calculate the normalised peak for each channel. This involves reading through the whole file which can be slow on large files.
func (f *File) CalcNormMaxAllChannels() (ret []float64, err error) {
	c := f.Format.Channels
	ret = make([]float64, c)
	e := C.sf_command(f.s, C.SFC_CALC_NORM_MAX_ALL_CHANNELS, unsafe.Pointer(&ret[0]), C.int(c*8))
	if e != 0 {
		err = sErrorType(e)
	}
	return
}

//Retrieve the peak value for the file as stored in the file header.
func (f *File) GetSignalMax() (ret float64, ok bool) {
	r := C.sf_command(f.s, C.SFC_GET_SIGNAL_MAX, unsafe.Pointer(&ret), 8)
	if r == C.SF_TRUE {
		ok = true
	}
	return
}

//Retrieve the peak value for the file as stored in the file header.
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
func (f *File) SetFloatNormalization(norm bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_NORM_FLOAT, norm)
}

/*This command only affects data read from or written to using ReadItems, ReadFrames, WriteItems, or WriteFrames with slices of float64.

For read operations setting normalisation to true means that the data from all subsequent reads will be be normalised to the range [-1.0, 1.0].

For write operations, setting normalisation to true means than all data supplied to the float write functions should be in the range [-1.0, 1.0] and will be scaled for the file format as necessary.

For both cases, setting normalisation to false means that no scaling will take place.

Returns the previous normalization setting. */
func (f *File) SetDoubleNormalization(norm bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_NORM_DOUBLE, norm)
}

// Returns the current float32 normalization mode.
func (f *File) GetFloatNormalization() bool {
	return f.genericBoolBoolCmd(C.SFC_GET_NORM_FLOAT, false)
}

// Returns the current float64 normalization mode.
func (f *File) GetDoubleNormalization() bool {
	return f.genericBoolBoolCmd(C.SFC_GET_NORM_DOUBLE, false)
}

//Set/clear the scale factor when integer (short/int) data is read from a file containing floating point data.
func (f *File) SetFloatIntScaleRead(scale bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_SCALE_FLOAT_INT_READ, scale)
}

//Set/clear the scale factor when integer (short/int) data is written to a file as floating point data.
func (f *File) SetIntFloatScaleWrite(scale bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_SCALE_INT_FLOAT_WRITE, scale)
}

//Retrieve the number of simple formats supported by libsndfile.
func GetSimpleFormatCount() int {
	var o C.int
	C.sf_command(nil, C.SFC_GET_SIMPLE_FORMAT_COUNT, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o)))
	return int(o)
}

//Retrieve information about a simple format.
//The value of the format argument should be the format number (ie 0 <= format <= count value obtained using GetSimpleFormatCount()).
// The returned format argument is suitable for use in sndfile.Open()
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
func GetMajorFormatCount() int {
	var o C.int
	C.sf_command(nil, C.SFC_GET_FORMAT_MAJOR_COUNT, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o)))
	return int(o)
}

//Retrieve information about a major format type
//For a more comprehensive example, see the program list_formats.c in the examples/ directory of the libsndfile source code distribution.
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
func GetSubFormatCount() int {
	var o C.int
	C.sf_command(nil, C.SFC_GET_FORMAT_SUBTYPE_COUNT, unsafe.Pointer(&o), C.int(unsafe.Sizeof(o)))
	return int(o)
}

//Enumerate the subtypes (this function does not translate a subtype into a string describing that subtype). A typical use case might be retrieving a string description of all subtypes so that a dialog box can be filled in.
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
func (f *File) SetAddPeakChunk(set bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_ADD_PEAK_CHUNK, set)
}

//The header of an audio file is normally written by libsndfile when the file is closed using sf_close().

//There are however situations where large files are being generated and it would be nice to have valid data in the header before the file is complete. Using this command will update the file header to reflect the amount of data written to the file so far. Other programs opening the file for read (before any more data is written) will then read a valid sound file header.
func (f *File) UpdateHeaderNow() {
	C.sf_command(f.s, C.SFC_UPDATE_HEADER_NOW, nil, 0)
}

//Similar to SFC_UPDATE_HEADER_NOW but updates the header at the end of every call to the sf_write* functions.
func (f *File) SetUpdateHeaderAuto(set bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_UPDATE_HEADER_AUTO, set)
}

// Truncates a file to /count/ frames.  After this command, both the read and the write pointer will be at the new end of the file. This command will fail (returning non-zero) if the requested truncate position is beyond the end of the file.
func (f *File) Truncate(count int64) (err error) {
	r := C.sf_command(f.s, C.SFC_FILE_TRUNCATE, unsafe.Pointer(&count), 8)

	if r != 0 {
		err = errors.New(C.GoString(C.sf_strerror(f.s)))
	}
	return
}

func (f *File) genericBoolBoolCmd(cmd C.int, i bool) bool {
	ib := C.SF_FALSE
	if i {
		ib = C.SF_TRUE
	}

	n := C.sf_command(f.s, cmd, nil, C.int(ib))
	return (n == C.SF_TRUE)
}

//Change the data start offset for files opened up as SF_FORMAT_RAW. libsndfile implements this but it appears to not do anything useful that you can't accomplish with seek, so consider this deprecated.
func (f *File) SetRawStartOffset(count int64) (err error) {
	r := C.sf_command(f.s, C.SFC_SET_RAW_START_OFFSET, unsafe.Pointer(&count), 8)

	if r != 0 {
		err = errors.New(C.GoString(C.sf_strerror(f.s)))
	}
	return
}

//Turn on/off automatic clipping when doing floating point to integer conversion.
func (f *File) SetClipping(clip bool) bool {
	return f.genericBoolBoolCmd(C.SFC_SET_CLIPPING, clip)
}

//Is automatic clipping when doing floating point to integer conversion on?
func (f *File) GetClipping(clip bool) bool {
	return f.genericBoolBoolCmd(C.SFC_GET_CLIPPING, false)
}

//Get the file offset and file length of a file enbedded within another larger file.
//The value of the offset return value will be the offsets in bytes from the start of the outer file to the start of the embedded audio file.
//The value of the length return value will be the length in bytes of the embedded file.
// Untested.
func (f *File) GetEmbeddedFileInfo() (offset, length int64, err error) {
	var s C.SF_EMBED_FILE_INFO
	r := C.sf_command(f.s, C.SFC_GET_EMBED_FILE_INFO, unsafe.Pointer(&s), C.int(unsafe.Sizeof(s)))
	if r != 0 {
		err = errors.New(C.GoString(C.sf_strerror(f.s)))
	}
	offset = int64(s.offset)
	length = int64(s.length)
	return
}

const AmbisonicNone int = int(C.SF_AMBISONIC_NONE)
const AmbisonicBFormat int = int(C.SF_AMBISONIC_B_FORMAT)

//Test if the current file has the GUID of a WAVEX file for any of the Ambisonic formats.
// returns AmbisonicNone or AmbisonicBFormat, or zero if the file format does not support Ambisonic formats
func (f *File) WavexGetAmbisonic() int {
	return int(C.sf_command(f.s, C.SFC_WAVEX_GET_AMBISONIC, nil, 0))
}

//Set the GUID of a new WAVEX file to indicate an Ambisonics format.
// returns format that was just set, or zero if the file format does not support Ambisonic formats
func (f *File) WavexSetAmbisonic(ambi int) int {
	return int(C.sf_command(f.s, C.SFC_WAVEX_SET_AMBISONIC, nil, C.int(ambi)))
}

//Set the the Variable Bit Rate encoding quality. The encoding quality value should be between 0.0 (lowest quality) and 1.0 (highest quality). Untested.
func (f *File) SetVbrQuality(q float64) (err error) {
	r := C.sf_command(f.s, C.SFC_SET_VBR_ENCODING_QUALITY, unsafe.Pointer(&q), 8)
	if r != 0 {
		err = errors.New(C.GoString(C.sf_strerror(f.s)))
	}
	return
}

//Determine if raw data read using sf_read_raw needs to be end swapped on the host CPU.

//For instance, will return true on when reading WAV containing SF_FORMAT_PCM_16 data on a big endian machine and false on a little endian machine.
func (f *File) RawNeedsEndianSwap() bool {
	return f.genericBoolBoolCmd(C.SFC_RAW_DATA_NEEDS_ENDSWAP, false)
}

type BroadcastInfo struct {
	Description          string
	Originator           string
	Originator_reference string
	Origination_date     string
	Origination_time     string
	Time_reference_low   uint32
	Time_reference_high  uint32
	Version              uint16
	Umid                 string
	Coding_history       []int8
}

func goStringFromArr(c []C.char) string {
	s := make([]byte, len(c))
	for i, r := range c {
		s[i] = byte(r)
	}
	return string(s)
}

func trim(in string) string {
	i := strings.IndexRune(in, 0)
	if i > -1 {
		return in[0:i]
	} else {
		return in
	}
	return "shelzbut" // shut up compiler
}

func broadcastFromC(c *C.SF_BROADCAST_INFO) *BroadcastInfo {
	bi := new(BroadcastInfo)
	bi.Description = trim(C.GoStringN(&c.description[0], C.int(len(c.description[:]))))
	bi.Originator = trim(C.GoStringN(&c.originator[0], C.int(len(c.originator[:]))))
	bi.Originator_reference = trim(C.GoStringN(&c.originator_reference[0], C.int(len(c.originator_reference[:]))))
	bi.Origination_date = trim(C.GoStringN(&c.origination_date[0], C.int(len(c.origination_date[:]))))
	bi.Origination_time = trim(C.GoStringN(&c.origination_time[0], C.int(len(c.origination_time[:]))))
	bi.Time_reference_low = uint32(c.time_reference_low)
	bi.Time_reference_high = uint32(c.time_reference_high)
	bi.Version = uint16(c.version)
	bi.Umid = trim(C.GoStringN(&c.umid[0], C.int(len(c.umid[:]))))
	bi.Coding_history = make([]int8, c.coding_history_size, 256)
	for i, r := range c.coding_history {
		if i >= int(c.coding_history_size) {
			break
		}
		bi.Coding_history = append(bi.Coding_history, int8(r))
	}
	return bi
}

// Retrieve the Broadcast Extension Chunk from WAV (and related) files.
func (f *File) GetBroadcastInfo() (bi *BroadcastInfo, ok bool) {
	bic := new(C.SF_BROADCAST_INFO)

	r := C.sf_command(f.s, C.SFC_GET_BROADCAST_INFO, unsafe.Pointer(bic), C.int(unsafe.Sizeof(*bic)))
	if r == C.SF_TRUE {
		bi = broadcastFromC(bic)
		ok = true
	}
	return
}

func arrFromGoString(arr []C.char, src string) {
	for i, r := range src {
		if i >= len(arr) {
			break
		}
		arr[i] = C.char(r)
	}
}

func cFromBroadcast(bi *BroadcastInfo) (c *C.SF_BROADCAST_INFO) {
	c = new(C.SF_BROADCAST_INFO)
	arrFromGoString(c.description[:], bi.Description)
	arrFromGoString(c.originator[:], bi.Originator)
	arrFromGoString(c.originator_reference[:], bi.Originator_reference)
	arrFromGoString(c.origination_date[:], bi.Origination_date)
	arrFromGoString(c.origination_time[:], bi.Origination_time)
	c.time_reference_low = C.uint(bi.Time_reference_low)
	c.time_reference_high = C.uint(bi.Time_reference_high)
	c.version = C.short(bi.Version)
	arrFromGoString(c.umid[:], bi.Umid)
	ch := bi.Coding_history
	if len(bi.Coding_history) > 256 {
		ch = bi.Coding_history[0:256]
	}
	c.coding_history_size = C.uint(len(ch))
	for i, r := range ch {
		c.coding_history[i] = C.char(r)
	}
	return c
}

// Set the Broadcast Extension Chunk from WAV (and related) files.
func (f *File) SetBroadcastInfo(bi *BroadcastInfo) (err error) {
	c := cFromBroadcast(bi)
	r := C.sf_command(f.s, C.SFC_SET_BROADCAST_INFO, unsafe.Pointer(c), C.int(unsafe.Sizeof(*c)))
	if r == C.SF_FALSE {
		err = errors.New(C.GoString(C.sf_strerror(f.s)))
	}
	return
}

type LoopMode int

const (
	None        LoopMode = C.SF_LOOP_NONE
	Forward              = C.SF_LOOP_FORWARD
	Backward             = C.SF_LOOP_BACKWARD
	Alternating          = C.SF_LOOP_ALTERNATING
)

type LoopInfo struct {
	TimeSig struct {
		Numerator   int16 // any positive integer
		Denominator int16 // any positive power of 2
	}
	Mode    LoopMode
	Beats   int     // not amount of quarter notes. a full bar of 7/8 is 7 bears
	Bpm     float32 // Suggestion
	RootKey int     // MIDI Note
	Future  [6]int  // nuffink
}

//Retrieve loop information for file including time signature, length in beats and original MIDI base note

// Returns populated structure if file contains loop info, otherwise nil. Untested.
func (f *File) GetLoopInfo() (i *LoopInfo) {
	c := new(C.SF_LOOP_INFO)
	r := C.sf_command(f.s, C.SFC_GET_LOOP_INFO, unsafe.Pointer(c), C.int(unsafe.Sizeof(*c)))
	if r == C.SF_TRUE {
		i = new(LoopInfo)
		i.TimeSig.Numerator = int16(c.time_sig_num)
		i.TimeSig.Denominator = int16(c.time_sig_den)
		i.Mode = LoopMode(c.loop_mode)
		i.Beats = int(c.num_beats)
		i.Bpm = float32(c.bpm)
		i.RootKey = int(c.root_key)
		for index, value := range c.future {
			i.Future[index] = int(value)
		}
	}
	return
}

type Instrument struct {
	Gain             int
	Basenote, Detune int8
	Velocity         [2]int8 // low byte is index 0
	Key              [2]int8 // low byte is index 0
	LoopCount        int
	Loops            [16]struct {
		Mode              LoopMode
		Start, End, Count uint
	}
}

// Retrieve instrument information from file including MIDI base note, keyboard mapping and looping information (start/stop and mode).

// Return pointer to populated structure if the file header contains instrument information for the file. nil otherwise.
func (f *File) GetInstrument() (i *Instrument) {
	c := new(C.SF_INSTRUMENT)
	i = new(Instrument)
	r := C.sf_command(f.s, C.SFC_GET_INSTRUMENT, unsafe.Pointer(c), C.int(unsafe.Sizeof(*c)))
	if r == C.SF_TRUE {
		i.Gain = int(c.gain)
		i.Basenote = int8(c.basenote)
		i.Detune = int8(c.detune)
		i.Velocity[0] = int8(c.velocity_lo)
		i.Velocity[1] = int8(c.velocity_hi)
		i.Key[0] = int8(c.key_lo)
		i.Key[1] = int8(c.key_hi)
		i.LoopCount = int(c.loop_count)
		for index, loop := range c.loops {
			i.Loops[index].Mode = LoopMode(loop.mode)
			i.Loops[index].Start = uint(loop.start)
			i.Loops[index].End = uint(loop.end)
			i.Loops[index].Count = uint(loop.count)
		}
	}
	return
}

// Set instrument information from file including MIDI base note, keyboard mapping and looping information (start/stop and mode).

// Return true if the file header contains instrument information for the file. false otherwise.
func (f *File) SetInstrument(i *Instrument) bool {
	c := new(C.SF_INSTRUMENT)
	c.gain = C.int(i.Gain)
	c.basenote = C.char(i.Basenote)
	c.detune = C.char(i.Detune)
	c.velocity_lo = C.char(i.Velocity[0])
	c.velocity_hi = C.char(i.Velocity[1])
	c.key_lo = C.char(i.Key[0])
	c.key_hi = C.char(i.Key[1])
	c.loop_count = C.int(i.LoopCount)
	var index int
	for ; index < i.LoopCount; index++ {
		c.loops[index].mode = C.int(i.Loops[index].Mode)
		c.loops[index].start = C.uint(i.Loops[index].Start)
		c.loops[index].end = C.uint(i.Loops[index].End)
		c.loops[index].count = C.uint(i.Loops[index].Count)
	}
	for ; index < 16; index++ {
		c.loops[index].mode = C.int(None)
		// why is this necessary? libsndfile doesn't check loopcount for AIFF
	}

	r := C.sf_command(f.s, C.SFC_SET_INSTRUMENT, unsafe.Pointer(c), C.int(unsafe.Sizeof(*c)))
	return (r == C.SF_TRUE)
}

// This allows libsndfile experts to use the command interface for commands not currently supported. See http://www.mega-nerd.com/libsndfile/command.html
// The f argument may be nil in cases where the command does not require a SNDFILE argument.
// The method's cmd, data, and datasize arguments are used the same way as the correspondingly named arguments for sf_command
func GenericCmd(f *File, cmd C.int, data unsafe.Pointer, datasize int) int {
	var s *C.SNDFILE = nil
	if f != nil {
		s = f.s
	}
	return int(C.sf_command(s, cmd, data, C.int(datasize)))
}

const (
	ChannelMapInvalid = C.SF_CHANNEL_MAP_INVALID
	ChannelMapMono = C.SF_CHANNEL_MAP_MONO
	ChannelMapLeft = C.SF_CHANNEL_MAP_LEFT /* Apple calls this 'Left' */
	ChannelMapRight = C.SF_CHANNEL_MAP_RIGHT /* Apple calls this 'Right' */
	ChannelMapCenter = C.SF_CHANNEL_MAP_CENTER /* Apple calls this 'Center' */
	ChannelMapFrontLeft = C.SF_CHANNEL_MAP_FRONT_LEFT
	ChannelMapFrontRight = C.SF_CHANNEL_MAP_FRONT_RIGHT
	ChannelMapFrontCenter = C.SF_CHANNEL_MAP_FRONT_CENTER
	ChannelMapRearCenter = C.SF_CHANNEL_MAP_REAR_CENTER /* Apple calls this 'Center Surround' Msft calls this 'Back Center' */
	ChannelMapRearLeft = C.SF_CHANNEL_MAP_REAR_LEFT /* Apple calls this 'Left Surround' Msft calls this 'Back Left' */
	ChannelMapRearRight = C.SF_CHANNEL_MAP_REAR_RIGHT /* Apple calls this 'Right Surround' Msft calls this 'Back Right' */
	ChannelMapLfe = C.SF_CHANNEL_MAP_LFE /* Apple calls this 'LFEScreen' Msft calls this 'Low Frequency'  */
	ChannelMapFrontLeftOfCenter = C.SF_CHANNEL_MAP_FRONT_LEFT_OF_CENTER /* Apple calls this 'Left Center' */
	ChannelMapFrontRightOfCenter = C.SF_CHANNEL_MAP_FRONT_RIGHT_OF_CENTER /* Apple calls this 'Right Center */
	ChannelMapSideLeft = C.SF_CHANNEL_MAP_SIDE_LEFT /* Apple calls this 'Left Surround Direct' */
	ChannelMapSideRight = C.SF_CHANNEL_MAP_SIDE_RIGHT /* Apple calls this 'Right Surround Direct' */
	ChannelMapTopCenter = C.SF_CHANNEL_MAP_TOP_CENTER /* Apple calls this 'Top Center Surround' */
	ChannelMapTopFrontLeft = C.SF_CHANNEL_MAP_TOP_FRONT_LEFT /* Apple calls this 'Vertical Height Left' */
	ChannelMapTopFrontRight = C.SF_CHANNEL_MAP_TOP_FRONT_RIGHT /* Apple calls this 'Vertical Height Right' */
	ChannelMapTopFrontCenter = C.SF_CHANNEL_MAP_TOP_FRONT_CENTER /* Apple calls this 'Vertical Height Center' */
	ChannelMapTopRearLeft = C.SF_CHANNEL_MAP_TOP_REAR_LEFT /* Apple and MS call this 'Top Back Left' */
	ChannelMapTopRearRight = C.SF_CHANNEL_MAP_TOP_REAR_RIGHT /* Apple and MS call this 'Top Back Right' */
	ChannelMapTopRearCenter = C.SF_CHANNEL_MAP_TOP_REAR_CENTER /* Apple and MS call this 'Top Back Center' */
	ChannelMapAmbisonicBW = C.SF_CHANNEL_MAP_AMBISONIC_B_W
	ChannelMapAmbisonicBX = C.SF_CHANNEL_MAP_AMBISONIC_B_X
	ChannelMapAmbisonicBY = C.SF_CHANNEL_MAP_AMBISONIC_B_Y
	ChannelMapAmbisonicBZ = C.SF_CHANNEL_MAP_AMBISONIC_B_Z
	ChannelMapMax = C.SF_CHANNEL_MAP_MAX
)

// Returns a slice full of integers detailing the position of each channel in the file. err will be non-nil on an actual error
func (f *File) GetChannelMapInfo() (channels []int32, err error) {
	channels = make([]int32, f.Format.Channels)
	r := GenericCmd(f, C.SFC_GET_CHANNEL_MAP_INFO, unsafe.Pointer(&channels[0]), len(channels)*4)
	if r == C.SF_FALSE {
		err = errors.New(C.GoString(C.sf_strerror(f.s)))
	}
	return
}

func (f *File) SetChannelMapInfo(channels []int32) (err error) {
	if int32(len(channels)) != f.Format.Channels {
		err = errors.New("channel map passed in didn't match file channel count " + string(len(channels)) + " != " + string(f.Format.Channels))
	}
	r := GenericCmd(f, C.SFC_SET_CHANNEL_MAP_INFO, unsafe.Pointer(&channels[0]), len(channels)*4)
	if r == C.SF_FALSE {
		err = errors.New(C.GoString(C.sf_strerror(f.s)))
	}
	return
}
