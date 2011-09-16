package sndfile

// #include <sndfile.h>
// #include "virtual.h"
import "C"
import "unsafe"
import "os"


type VIO_get_filelen func(interface{}) int64
type VIO_seek func(int64, Whence, interface{}) int64
type VIO_read func([]byte, interface{}) int64
type VIO_write func([]byte, interface{}) int64
type VIO_tell func(interface{}) int64

// Opens a soundfile from a virtual file I/O context which is provided by the caller. This is usually used to interface libsndfile to a stream or buffer based system. Apart from the c and user_data parameters this function behaves like sf_open.
// THIS PART OF THE PACKAGE IS EXPERIMENTAL. Don't use it yet.
// needs test. lots of them
func OpenVirtual(v VirtualIo, mode Mode, info Info, user_data interface{}) (f *File, err os.Error) {
	c := C.virtualio()
	var vp virtualIo
	vp.v = &v
	vp.c = &c
	f = new(File)
	f.s = C.sf_open_virtual(&c, C.int(mode), info.toCinfo(), unsafe.Pointer(&vp))
	if f.s != nil {
		f.virtual = &vp
	} else {
		err = sErrorType(C.sf_error(nil))
	}
	return
}

// You must provide the following:
//UserData is the virtual file context. It is opaque to this layer.
//GetLength returns the length of the virtual file in BYTES
//Seek - The virtual file context must seek to offset using the seek mode provided by whence which is one of
//Read - The virtual file context must copy ("read") "count" bytes into the buffer provided by ptr and return the count of actually copied bytes. (only when file is opened in Read or ReadWrite mode)
//Write - The virtual file context must process "count" bytes stored in the buffer passed with ptr and return the count of actually processed bytes.
//Tell - Return the current position of the virtual file context.
type VirtualIo struct {
	GetLength VIO_get_filelen 
	Seek VIO_seek 
	Read VIO_read 
	Write VIO_write
	Tell VIO_tell
	UserData interface{}
}

type virtualIo struct {
	v *VirtualIo
	c *C.SF_VIRTUAL_IO
}

//export gsfLen
func gsfLen (user_data unsafe.Pointer) int64 {
	l := (*VirtualIo)(user_data)
	return l.GetLength(l.UserData)
}

//export gsfSeek
func gsfSeek (i int64, w Whence, user_data unsafe.Pointer) int64 {
	l := (*VirtualIo)(user_data)
	return l.Seek(i, w, l.UserData)
}

//export gsfRead
func gsfRead (ptr unsafe.Pointer, i int64, user_data unsafe.Pointer) int64 {
	l := (*VirtualIo)(user_data)
	b := (*[1<<30]byte)(ptr)[0:i]
	return l.Read(b, l.UserData)
}

//export gsfWrite
func gsfWrite(ptr unsafe.Pointer, i int64, user_data unsafe.Pointer) int64 {
	l := (*VirtualIo)(user_data)
	b := (*[1<<30]byte)(ptr)[0:i]
	return l.Write(b, l.UserData)
}

//export gsfTell
func gsfTell(user_data unsafe.Pointer) int64 {
	l := (*VirtualIo)(user_data)
	return l.Tell(l.UserData)
}

