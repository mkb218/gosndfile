package sndfile

// #include <sndfile.h>
// #include "virtual.h"
import "C"
//import "fmt"
import "runtime"
import "unsafe"
import "errors"

type VIO_get_filelen func(interface{}) int64
type VIO_seek func(int64, Whence, interface{}) int64
type VIO_read func([]byte, interface{}) int64
type VIO_write func([]byte, interface{}) int64
type VIO_tell func(interface{}) int64

// Opens a soundfile from a virtual file I/O context which is provided by the caller. This is usually used to interface libsndfile to a stream or buffer based system. Apart from the c and user_data parameters this function behaves like sf_open.
// THIS PART OF THE PACKAGE IS EXPERIMENTAL. Don't use it yet.
func OpenVirtual(v VirtualIo, mode Mode, info *Info) (f *File, err error) {
	c := C.virtualio()
	var vp virtualIo
	vp.v = &v
	vp.c = c
	f = new(File)
	ci := info.toCinfo()
	f.s = C.sf_open_virtual(c, C.int(mode), ci, unsafe.Pointer(&vp))
	if f.s != nil {
		f.virtual = &vp
		f.Format = fromCinfo(ci)
		*info = f.Format
	} else {
		err = errors.New(C.GoString(C.sf_strerror(nil)))
	}
	runtime.SetFinalizer(f, (*File).Close)
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
	Seek      VIO_seek
	Read      VIO_read
	Write     VIO_write
	Tell      VIO_tell
	UserData  interface{}
}

type virtualIo struct {
	v *VirtualIo
	c *C.SF_VIRTUAL_IO
}

//export gsfLen
func gsfLen(user_data unsafe.Pointer) int64 {
	l := (*virtualIo)(user_data)
	return l.v.GetLength(l.v.UserData)
}

//export gsfSeek
func gsfSeek(i int64, w Whence, user_data unsafe.Pointer) int64 {
	if user_data == nil {
		panic("nil ud")
	}
	l := (*virtualIo)(user_data)
	return l.v.Seek(i, w, l.v.UserData)
}

//export gsfRead
func gsfRead(ptr unsafe.Pointer, i int64, user_data unsafe.Pointer) int64 {
	l := (*virtualIo)(user_data)
	b := (*[1 << 30]byte)(ptr)[0:i]
	return l.v.Read(b, l.v.UserData)
}

//export gsfWrite
func gsfWrite(ptr unsafe.Pointer, i int64, user_data unsafe.Pointer) int64 {
	l := (*virtualIo)(user_data)
	b := (*[1 << 30]byte)(ptr)[0:i]
	return l.v.Write(b, l.v.UserData)
}

//export gsfTell
func gsfTell(user_data unsafe.Pointer) int64 {
	l := (*virtualIo)(user_data)
	return l.v.Tell(l.v.UserData)
}
