package sndfile

import "testing"
import "fmt"
import "strings"
import "unsafe"

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
	fmt.Println("--- Supported simple formats")
	for i := 0; i < simpleformats; i++ {
		f, name, ext, ok := GetSimpleFormat(i)
		fmt.Printf("%08x %s %s\n", f, name, ext)
		if !ok {
			t.Error("error from GetSimpleFormat()")
		}
	}
	
	fmt.Println("--- Supported formats")
	// following is straight from examples in libsndfile distribution
	majorcount := GetMajorFormatCount()
	subcount := GetSubFormatCount()
	for m := 0; m < majorcount; m++ {
		f, name, ext, ok := GetMajorFormatInfo(m)
		if ok {
			fmt.Printf("--- MAJOR 0x%08x %v Extension: .%v\n", f, name, ext)
			for s := 0; s < subcount; s++ {
				t, sname, sok := GetSubFormatInfo(s)
				var i Info
				i.Channels = 1
				i.Format = Format(f|t)
				if sok && FormatCheck(i) {
					fmt.Printf("   0x%08x %v %v\n", f|t, name, sname)
				} else {
//					fmt.Printf("no format pair 0x%x\n", f|t)
				}
			}
		} else {
			fmt.Printf("no format for number %v\n", m)
		}
	}
}

func isLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return (b == 0x04)
}

func TestRawSwap(t *testing.T) {
	// set up file to be checked
	i := &Info{0, 44100, 1, SF_FORMAT_WAV|SF_FORMAT_PCM_16,0,0}
	f, err := Open("leout.wav", Write, i)
	if err != nil {
		t.Fatalf("couldn't open file for writing: %v", err)
	}
	if isLittleEndian() && f.RawNeedsEndianSwap() {
		t.Errorf("little endian file and little endian cpu shuld not report needing swap but does!")
	} else if !isLittleEndian() && !f.RawNeedsEndianSwap() {
		t.Errorf("little endian file and big endian machine should report needing swap, but doesn't!")
	}
}

func TestGenericCmd(t *testing.T) {
	i := GenericCmd(nil, 0x1000, nil, 0)
	c := make([]byte, i)
	GenericCmd(nil, 0x1000, unsafe.Pointer(&c[0]), i)
	if !strings.HasPrefix(string(c), "libsndfile") {
		t.Errorf("version string \"%s\" had unexpected prefix", string(c))
	}	
}