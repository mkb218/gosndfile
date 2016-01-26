package sndfile

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

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
	//	fmt.Printf("getUd %v\n", ud.f)
	return ud
}

func testGetLength(i interface{}) int64 {
	//	fmt.Println("gogetlength")
	ud := getUd(i)
	s, err := ud.f.Stat()
	if err != nil {
		// is this even possible?
		ud.t.Errorf("couldn't get length for some reason %s", err)
	}
	return int64(s.Size())
}

func testSeek(offset int64, whence Whence, i interface{}) int64 {
	//	fmt.Printf("goseek %d %d\n", offset, whence)
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
		ud.t.Errorf("Couldn't seek for some reason: %s", err)
	}
	return o
}

func testRead(buf []byte, i interface{}) int64 {
	//	fmt.Println("goread")
	ud := getUd(i)
	read, err := ud.f.Read(buf)
	if err != nil {
		ud.t.Errorf("couldn't read from file %s", err)
	}
	return int64(read)
}

func testWrite(buf []byte, i interface{}) int64 {
	//	fmt.Println("gowrite")
	ud := getUd(i)
	wrote, err := ud.f.Write(buf)
	if err != nil {
		ud.t.Errorf("couldn't write to file %v %d %s", ud.f, wrote, err)
	}
	return int64(wrote)
}

func testTell(i interface{}) int64 {
	//	fmt.Println("gotell")
	ud := getUd(i)
	o, err := ud.f.Seek(0, os.SEEK_CUR)
	if err != nil {
		ud.t.Errorf("couldn't tell! %s", err)
	}
	return o
}

// test virtual i/o by mapping virtual i/o calls to Go i/o calls
func TestVirtualRead(t *testing.T) {
	f, err := os.Open("test/ok.aiff")
	if err != nil {
		t.Fatalf("couldn't open input file %s", err)
	}

	var vi VirtualIo
	vi.UserData = testUserData{f, t}
	vi.GetLength = testGetLength
	vi.Seek = testSeek
	vi.Read = testRead
	vi.Write = testWrite
	vi.Tell = testTell

	var i Info
	vf, err := OpenVirtual(vi, Read, &i)
	if err != nil {
		t.Fatalf("error from OpenVirtual %v", err)
	}
	if !reflect.DeepEqual(i, goldenInfo()) {
		t.Errorf("info struct not as expected! %v vs. golden %v", i, goldenInfo())
	}
	off, err := vf.Seek(0, Set)
	if off != 0 || err != nil {
		t.Errorf("Seek had wrong result %v (expected 0) %v", off, err)
	}
}

// test virtual i/o by mapping virtual i/o calls to Go i/o calls
func TestVirtualWrite(t *testing.T) {
	f, err := os.Create("test/funky2.aiff")
	if err != nil {
		t.Fatalf("couldn't open input file test/funky2.aiff %s", err)
	}

	var vi VirtualIo
	vi.UserData = testUserData{f, t}
	vi.GetLength = testGetLength
	vi.Seek = testSeek
	vi.Read = testRead
	vi.Write = testWrite
	vi.Tell = testTell

	var i Info
	i.Samplerate = 44100
	i.Channels = 2
	i.Format = SF_FORMAT_AIFF | SF_FORMAT_FLOAT
	vf, err := OpenVirtual(vi, Write, &i)
	if err != nil {
		t.Fatalf("error from OpenVirtual %v", err)
	}
	off, err := vf.Seek(0, Set)
	if off != 0 || err != nil {
		t.Errorf("Seek had wrong result %v (expected 0) %v", off, err)
	}

	out := []float32{
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.5, 0.0,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5,
		0.0, 0.5}
	out = append(out, out...)
	out = append(out, out...)
	out = append(out, out...)
	out = append(out, out...)
	out = append(out, out...)
	written, err := vf.WriteItems(out)
	if written != int64(len(out)) {
		t.Errorf("unexpected written item count %d not %d\n", written, len(out))
	}
	err = vf.Close()
	if err != nil {
		t.Errorf("virtual close failed %s\n", err)
	}
	err = f.Close()
	if err != nil {
		t.Errorf("close failed %s\n", err)
	}
	var ri Info
	_, err = Open("test/funky.aiff", Read, &ri)
	if err != nil {
		t.Fatalf("couldn't open input file %s", err)
	}

	//	fmt.Println(ri)
	//	fmt.Println(rf)

	if ri.Frames != int64(len(out)/2) {
		t.Errorf("length in samples not as expected! %d vs. expected %d", ri.Frames, len(out)/2)
	}
}
