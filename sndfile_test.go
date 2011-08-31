package sndfile

import (
//	"fmt"
	"reflect"
	"testing"
)

func goldenInfo() (i Info) {
	i.Frames = 24036
	i.Samplerate = 8012
	i.Channels = 1
	i.Format = 131074
	i.Sections = 1
	i.Seekable = 1
	return
}

func goldenShortInput() []int16 {
	return []int16{0,0,0,0,0}
}

func TestReadShortItems(t *testing.T) {
	var i Info
	f, e := Open("ok.aiff", Read, &i)
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