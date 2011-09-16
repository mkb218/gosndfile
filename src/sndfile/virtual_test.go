package sndfile

import "fmt"
import "testing"
import "os"
import "reflect"

type testUserData struct {
	f *os.File
	t *testing.T
}

func getUd(i interface{}) testUserData {
	ud, ok := i.(testUserData)
	if !ok {
		// i can't even guarantee that I can get a t out of it so just die
		fmt.Fprintf(os.Stderr, "userdata didn't contain a valid struct! %v\n", reflect.TypeOf(i))
		os.Exit(1)
	}
	return ud
}

func testGetLength(i interface{}) int64 {
	ud := getUd(i)
	s, err := ud.f.Stat()
	if err != nil {
		// is this even possible?
		ud.t.Errorf("couldn't get length for some reason %s", err.String())
	}
	return int64(s.Size)
}

func testSeek(offset int64, whence Whence, i interface{}) int64 {
	ud := getUd(i)
	var w int
	if whence == Set {
		w = os.SEEK_SET
	} else if whence == Current {
		w = os.SEEK_CUR
	} else if whence == End {
		w = os.SEEK_END
	}
	
	o, err := ud.f.Seek(offset, w)
	if err != nil {
		ud.t.Errorf("Couldn't seek for some reason: %s", err.String())
	}
	return o
}

func testRead(buf []byte, i interface{}) int64 {
	ud := getUd(i)
	read, err := ud.f.Read(buf)
	if err != nil {
		ud.t.Errorf("couldn't read from file %s", err.String())
	}
	return int64(read)
}

func testWrite(buf []byte, i interface{}) int64 {
	ud := getUd(i)
	wrote, err := ud.f.Write(buf)
	if err != nil {
		ud.t.Errorf("couldn't write to file %s", err.String())
	}
	return int64(wrote)
}

func testTell(i interface{}) int64 {
	ud := getUd(i)
	o, err := ud.f.Seek(0, os.SEEK_CUR)
	if err != nil {
		ud.t.Errorf("couldn't tell! %s", err.String())
	}
	return o
}
	

// test virtual i/o by mapping virtual i/o calls to Go i/o calls
func TestVirtual(t *testing.T) {
	f, err := os.Open("ok.aiff")
	if err != nil {
		t.Fatalf("couldn't open input file %s", err.String())
	}
	
	var vi VirtualIo
	vi.UserData = testUserData{f, t}
	vi.GetLength = testGetLength
	vi.Seek = testSeek
	vi.Read = testRead
	vi.Write = testWrite
	vi.Tell = testTell
	
	var i Info
	vf, err := OpenVirtual(vi, ReadWrite, &i)
	if err != nil {
		t.Fatalf("error from OpenVirtual %v", err)
	}
	off, err := vf.Seek(0, Set)
	if off != 0 || err != nil {
		t.Errorf("Seek had wrong result %v (expected 0) %v", off, err)
	}
}