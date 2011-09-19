package sndfile

import "testing"
import "fmt"
import "strings"

func TestGetLibVersion(t *testing.T) {
	s, _ := GetLibVersion()
	fmt.Println(s)
	if !strings.HasPrefix(s, "libsndfile") {
		t.Errorf("version string \"%s\" had unexpected prefix", s)
	}
}

func TestGetLogInfo(t *testing.T) {
	var i Info
	f, err := Open("fwerefrg", Read, &i)
//	fmt.Println(f)
//	fmt.Println(err)
	s, err := f.GetLogInfo()
	fmt.Println("TestGetLogInfo output: ", s)
	if err != nil {
		t.Error("TestGetLogInfo err: ", err)
	}
}



func TestFileCommands(t *testing.T) {
	var i Info
	f, err := Open("funky.aiff", ReadWrite, &i)
	if err != nil {
		t.Fatalf("open file failed %s", err)
	}
	
	max, err := f.CalcSignalMax()
	if err != nil {
		t.Fatalf("signal max failed %s", err)
	}
	
	fmt.Printf("max signal %f\n", max)

	max, err = f.CalcNormSignalMax()
	if err != nil {
		t.Fatalf("norm signal max failed %s", err)
	}
	
	fmt.Printf("norm max signal %f\n", max)

	maxarr, err := f.CalcMaxAllChannels()
	if err != nil {
		t.Fatalf("max all chans failed %s", err)
	}
	
	fmt.Printf("max all chans signal %v\n", maxarr)
	
	maxarr, err = f.CalcNormMaxAllChannels()
	if err != nil {
		t.Fatalf("max all chans failed %s", err)
	}
	
	fmt.Printf("norm max all chans signal %v\n", maxarr)


	max, ok := f.GetSignalMax()
	if !ok {
		t.Error("got unexpected failure from GetSignalMax with val ", max)
	}

	maxarr, ok = f.GetMaxAllChannels()
	if !ok {
		t.Error("got unexpected failure from GetMaxAllChannels with vals", maxarr)
	}
}

func TestFormats(t *testing.T) {
	simpleformats := GetSimpleFormatCount()
	for i := 0; i < simpleformats; i++ {
		f, name, ext, ok := GetSimpleFormat(i)
		fmt.Printf("%08x %s %s\n", f, name, ext)
		if !ok {
			t.Error("error from GetSimpleFormat()")
		}
	}
	
	// following is straight from examples in libsndfile distribution
	majorcount := GetMajorFormatCount()
	subcount := GetSubFormatCount()
	for m := 0; m < majorcount; m++ {
		f, name, ext, ok := GetMajorFormatInfo(m)
	}
}