package sndfile

// #include <sndfile.h>
// #include "virtual.h"
import "C"
import "unsafe"

type VIO_get_filelen func(interface{}) int64
type VIO_seek func(int64, Whence, interface{}) int64
type VIO_read func(unsafe.Pointer, int64, interface{}) int64
type VIO_write func(unsafe.Pointer, int64, interface{}) int64
type VIO_tell func(interface{}) int64

// note : sndfile doesn't copy the SF_VIRTUAL_INFO you give it. make sure the struct won't get eaten by GC

type lenstr struct {
	callback VIO_get_filelen
	user_data interface{}
}

//export gsfLen
func gsfLen (user_data unsafe.Pointer) int64 {
	l := (*lenstr)(user_data)
	return l.callback(l.user_data)
}

type seekstr struct {
	callback VIO_seek
	user_data interface{}
}

//export gsfSeek
func gsfSeek (i int64, w Whence, user_data unsafe.Pointer) int64 {
	l := (*seekstr)(user_data)
	return l.callback(i, w, l.user_data)
}

type readstr struct {
	callback VIO_read
	user_data interface{}
}

//export gsfRead
func gsfRead (ptr unsafe.Pointer, i int64, user_data unsafe.Pointer) int64 {
	l := (*readstr)(user_data)
	return l.callback(ptr, i, l.user_data)
}

type writestr struct {
	callback VIO_write
	user_data interface{}
}

//export gsfWrite
func gsfWrite(ptr unsafe.Pointer, i int64, user_data unsafe.Pointer) int64 {
	l := (*writestr)(user_data)
	return l.callback(ptr, i, l.user_data)
}

type tellstr struct {
	callback VIO_tell
	user_data interface{}
}

//export gsfTell
func gsfTell(user_data unsafe.Pointer) int64 {
	l := (*tellstr)(user_data)
	return l.callback(l.user_data)
}

