package sndfile

import (
	"os"
	"reflect"
	"testing"
)

func goldenInfo() (i Info) {
	i.Frames = 24036
	i.Samplerate = 8012
	i.Channels = 1
	i.Format = 131074
	i.Sections = 0
	i.Seekable = 0
	return
}

func goldenShortInput() []int16 {
	return make([]int16, 5)
}

func goldenShortFramesInput() []int16 {
	return make([]int16, 10)
}

func goldenShortFramesSeekInput() []int16 {
	return []int16{-4860, -5884, -5884, -6140, -5884, -5884, -5372, -4860, -3772, -1756}
}

func TestReadShortItems(t *testing.T) {
	var i Info
	f, e := Open("test/ok.aiff", Read, &i)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(i, goldenInfo()) {
		t.Errorf("info struct not as expected! %v vs. golden %v", i, goldenInfo())
	}
	buf := make([]int16, 5)
	r, e := f.ReadItems(buf)
	if r != 5 {
		t.Errorf("only read %d out of 5 items", r)
	}
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(buf, goldenShortInput()) {
		t.Errorf("data not as expected! %v vs golden %v", buf, goldenShortInput())
	}
	return
}

func TestReadShortFrames(t *testing.T) {
	var i Info
	f, e := Open("test/ok.aiff", Read, &i)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(i, goldenInfo()) {
		t.Errorf("info struct not as expected! %v vs. golden %v", i, goldenInfo())
	}
	buf := make([]int16, 10)
	r, e := f.ReadFrames(buf)
	if r != 10 {
		t.Errorf("only read %d out of 10 items", r)
	}
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(buf, goldenShortFramesInput()) {
		t.Errorf("data not as expected! %v vs golden %v", buf, goldenShortFramesInput())
	}
	return
}

func TestReadShortFramesSeek(t *testing.T) {
	var i Info
	f, e := Open("test/ok.aiff", Read, &i)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(i, goldenInfo()) {
		t.Errorf("info struct not as expected! %v vs. golden %v", i, goldenInfo())
	}

	_, e = f.Seek(i.Frames/2, Current)
	if e != nil {
		t.Fatal(e)
	}

	buf := make([]int16, 10)
	r, e := f.ReadFrames(buf)
	if r != 10 {
		t.Errorf("only read %d out of 10 items", r)
	}
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(buf, goldenShortFramesSeekInput()) {
		t.Errorf("data not as expected! %v vs golden %v", buf, goldenShortFramesSeekInput())
	}
	return
}

func goldenIntInput() []int32 {
	return make([]int32, 5)
}

func goldenIntFramesInput() []int32 {
	return make([]int32, 10)
}

func goldenIntFramesSeekInput() []int32 {
	// difference in magnitude is because libsndfile makes the MOST significant bit the
	return []int32{-318504960, -385613824, -385613824, -402391040, -385613824, -385613824, -352059392, -318504960, -247201792, -115081216}
}

func TestReadIntItems(t *testing.T) {
	var i Info
	f, e := Open("test/ok.aiff", Read, &i)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(i, goldenInfo()) {
		t.Errorf("info struct not as expected! %v vs. golden %v", i, goldenInfo())
	}
	buf := make([]int32, 5)
	r, e := f.ReadItems(buf)
	if r != 5 {
		t.Errorf("only read %d out of 5 items", r)
	}
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(buf, goldenIntInput()) {
		t.Errorf("data not as expected! %v vs golden %v", buf, goldenIntInput())
	}
	return
}

func TestReadIntFrames(t *testing.T) {
	var i Info
	f, e := Open("test/ok.aiff", Read, &i)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(i, goldenInfo()) {
		t.Errorf("info struct not as expected! %v vs. golden %v", i, goldenInfo())
	}
	buf := make([]int32, 10)
	r, e := f.ReadFrames(buf)
	if r != 10 {
		t.Errorf("only read %d out of 10 items", r)
	}
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(buf, goldenIntFramesInput()) {
		t.Errorf("data not as expected! %v vs golden %v", buf, goldenIntFramesInput())
	}
	return
}

func TestReadIntFramesSeek(t *testing.T) {
	var i Info
	f, e := Open("test/ok.aiff", Read, &i)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(i, goldenInfo()) {
		t.Errorf("info struct not as expected! %v vs. golden %v", i, goldenInfo())
	}

	_, e = f.Seek(i.Frames/2, Current)
	if e != nil {
		t.Fatal(e)
	}

	buf := make([]int32, 10)
	r, e := f.ReadFrames(buf)
	if r != 10 {
		t.Errorf("only read %d out of 10 items", r)
	}
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(buf, goldenIntFramesSeekInput()) {
		t.Errorf("data not as expected! %v vs golden %v", buf, goldenIntFramesSeekInput())
	}
	return
}

// openfd
func TestOpenFd(t *testing.T) {
	osf, err := os.Create("openfd")
	if err != nil {
		t.Fatal("err opening file for fd", err)
	}
	var i Info
	i.Channels = 1
	i.Samplerate = 44100
	i.Format = SF_FORMAT_AIFF | SF_FORMAT_PCM_16
	snf, err := OpenFd(osf.Fd(), Write, &i, true)
	if err != nil {
		t.Fatal("err opening fd", err)
	}
	n, err := snf.WriteItems([]int16{1, 2, 3, 4, 5})
	if n != 5 || err != nil {
		t.Error("bad write", n, err)
	}
	snf.Close()
	nf := os.NewFile(osf.Fd(), "openfd")
	_, err = nf.Write([]byte{1, 2, 3, 4, 5})
	if err == nil {
		t.Error("File must not have closed")
	}
}

// set/get string
func TestGetSetString(t *testing.T) {
	s := "TEST_STRING_HOOAH"
	var i Info
	i.Format = SF_FORMAT_WAV | SF_FORMAT_PCM_16
	i.Channels = 1
	i.Samplerate = 44100
	f, err := Open("getsetstring.wav", Write, &i)
	if err != nil {
		t.Fatal("couldn't open file to write", err)
	}
	f.SetString(s, First)
	f.WriteItems([]int16{1, 2, 3, 4, 5})
	f.Close()

	f, err = Open("getsetstring.wav", Read, &i)
	if err != nil {
		t.Fatal("couldn't open file to read", err)
	}
	si := f.GetString(First)
	if si != s {
		t.Error("wrong string came back!", s, "!=", si)
	}
}

func TestError(t *testing.T) {
	var i Info
	_, err := Open("nonexistentfile", Read, &i)
	t.Log(err)
}
