package sndfile

import "os"
import "math"
import "reflect"
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
	f, err := Open("test/funky.aiff", ReadWrite, &i)
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
	
	f.Close()
	
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
	f.Close()
}

func TestGenericCmd(t *testing.T) {
	i := GenericCmd(nil, 0x1000, nil, 0)
	c := make([]byte, i)
	GenericCmd(nil, 0x1000, unsafe.Pointer(&c[0]), i)
	if !strings.HasPrefix(string(c), "libsndfile") {
		t.Errorf("version string \"%s\" had unexpected prefix", string(c))
	}	
}

func TestTruncate(t *testing.T) {
	// first write 100 samples to a file
	var i Info
	i.Samplerate = 44100
	i.Channels = 1
	i.Format = SF_FORMAT_AIFF|SF_FORMAT_PCM_24
	os.Remove("truncout.aiff")
	f, err := Open("truncout.aiff", ReadWrite, &i)
	if err != nil {
		t.Fatalf("couldn't open file for output! %v", err)
	}
	
	var junk [100]int32
	written, err := f.WriteItems(junk[0:100])
	if written != 100 {
		t.Errorf("wrong written count %d", written)
	}
	
	f.WriteSync()
	
	f.Truncate(20)

	f.WriteSync()
	
	seek, err := f.Seek(0, Current)
	if seek != 20 {
		t.Errorf("wrong seek %v", seek)
	}
	if err != nil {
		t.Errorf("error! %v", err)
	}
	f.Close()
}

func TestMax(t *testing.T) {
	// open file with no peak chunk
	var i Info
	i.Samplerate = 44100
	i.Channels = 4
	i.Format = SF_FORMAT_AIFF|SF_FORMAT_PCM_16
	
	f, err := Open("addpeakchunk1.aiff", Write, &i)
	if err != nil {
		t.Fatalf("couldn't open file %v", err)
	}
	
	f.SetAddPeakChunk(false)
	_, err = f.WriteItems([]int16{1,2,3,4,5,6,7,8,1,2,3,4,5,6,7,8,1,2,3,4,5,6,7,8,1,2,3,4,5,6,7,8})
	if err != nil {
		t.Error("write err:",err)
	}
	f.Close()

	f, err = Open("addpeakchunk1.aiff", Read, &i)
	if err != nil {
		t.Fatalf("couldn't open file %v", err)
	}
	
	// calc signals
	r, err := f.CalcSignalMax()
	if err != nil {
		t.Fatal("couldn't calculate signal max",err)
	}
	if r != 8.0 {
		t.Errorf("Signal max was %v\n", r)
	}
	
	r, err = f.CalcNormSignalMax()
	if err != nil {
		t.Fatal("couldn't calculate signal max",err)
	}
	if r != float64(8)/float64(0x8000) {
		t.Errorf("Signal max was %v not %v\n", r, float64(8)/float64(0x7fff))
	}
	
	ra, err := f.CalcMaxAllChannels()
	if err != nil {
		t.Fatal("couldn't calculate signal max",err)
	}
	if !reflect.DeepEqual(ra,[]float64{5.0,6.0,7.0,8.0}) {
		t.Errorf("Signal max was %v\n", ra)
	}
	
	ra, err = f.CalcNormMaxAllChannels()
	if err != nil {
		t.Fatal("couldn't calculate signal max",err)
	}
	if !reflect.DeepEqual(ra,[]float64{5.0/float64(0x8000),6.0/float64(0x8000),7.0/float64(0x8000),8.0/float64(0x8000)}) {
		t.Errorf("Signal max was %v\n", ra)
	}
	
	// make sure peak chunk returns false
	
	_, ok := f.GetSignalMax()
	if ok {
		t.Error("expected no peak chunk in file")
	}
	
	_, ok = f.GetMaxAllChannels()
	if ok {
		t.Error("expected no peak chunk in file")
	}

	f.Close()
	
	i.Format = SF_FORMAT_AIFF|SF_FORMAT_DOUBLE
	// repeat for peak chunk, making sure peak chunk returns same value
	f, err = Open("addpeakchunk1.aiff", Write, &i)
	if err != nil {
		t.Fatalf("couldn't open file %v", err)
	}
	
	f.SetAddPeakChunk(true)
	_, err = f.WriteItems([]int16{1,2,3,4,5,6,7,8,1,2,3,4,5,6,7,8,1,2,3,4,5,6,7,8,1,2,3,4,5,6,7,8})
	if err != nil {
		t.Error("write err:",err)
	}
	f.Close()

	f, err = Open("addpeakchunk1.aiff", Read, &i)
	if err != nil {
		t.Fatalf("couldn't open file %v", err)
	}
	
	// make sure peak chunk returns true
	
	r, ok = f.GetSignalMax()
	if !ok {
		t.Error("expected peak chunk in file")
	}
	if r != 8.0 { 
		t.Errorf("unexpected peak value %v", r)
	}	
	ra, ok = f.GetMaxAllChannels()
	if !ok {
		t.Error("expected peak chunk in file")
	}
	if !reflect.DeepEqual(ra,[]float64{5.0,6.0,7.0,8.0}) {
		t.Errorf("Signal max was %v\n", ra)
	}

}

func TestNormalization(t *testing.T) {
	var i Info
	i.Format = SF_FORMAT_AIFF|SF_FORMAT_PCM_16
	i.Channels = 1
	i.Samplerate = 44100
	f, err := Open("norm.aiff", Write, &i)
	if err != nil {
		t.Fatal("couldn't open write file", err)
	}
	
	// set write normalization on
	f.SetDoubleNormalization(true)
	f.SetFloatNormalization(true)
	_, err = f.WriteItems([]float64{-1,0,1,0,-1,0,1,0,-1,0,1,0,-1,0,1,0})
	if err != nil {
		t.Fatal("couldn't write to file", err)
	}
	_, err = f.WriteItems([]float32{-1,0,1,0,-1,0,1,0,-1,0,1,0,-1,0,1,0})
	if err != nil {
		t.Fatal("couldn't write to file", err)
	}

	// set write normalization off
	f.SetDoubleNormalization(false)
	f.SetFloatNormalization(false)
	_, err = f.WriteItems([]float64{-1,2,-3,4,-5,6,-7,8,-9,10,-11,12,-13,14,-15,16})
	if err != nil {
		t.Fatal("couldn't write to file", err)
	}
	_, err = f.WriteItems([]float32{-1,2,-3,4,-5,6,-7,8,-9,10,-11,12,-13,14,-15,16})
	if err != nil {
		t.Fatal("couldn't write to file", err)
	}
	
	f.Close()
	
	f, err = Open("norm.aiff", Read, &i)
	f.SetDoubleNormalization(false)
	f.SetFloatNormalization(false)
	f32 := make([]float32, 16)
	f64 := make([]float64, 16)
	
	f.ReadItems(f64)
	if !reflect.DeepEqual([]float64{-32767,0,32767,0,-32767,0,32767,0,-32767,0,32767,0,-32767,0,32767,0}, f64) {
		t.Errorf("read badness %v", f64)
	}
	f.ReadItems(f32)
	if !reflect.DeepEqual([]float32{-32767,0,32767,0,-32767,0,32767,0,-32767,0,32767,0,-32767,0,32767,0}, f32) {
		t.Errorf("read badness %v", f32)
	}
	f.ReadItems(f64)
	if !reflect.DeepEqual([]float64{-1,2,-3,4,-5,6,-7,8,-9,10,-11,12,-13,14,-15,16}, f64) {
		t.Errorf("read badness %v", f64)
	}
	f.ReadItems(f32)
	if !reflect.DeepEqual([]float32{-1,2,-3,4,-5,6,-7,8,-9,10,-11,12,-13,14,-15,16}, f32) {
		t.Errorf("read badness %v", f32)
	}
	
	f.Seek(0, Set)
	f.SetDoubleNormalization(true)
	f.SetFloatNormalization(true)
	
	n, err := f.ReadItems(f64)
	if err != nil || n != 16 {
		t.Fatal("bad read", err, n)
	}
	if !floatEqual([]float64{-float64(32767)/float64(32768),0,float64(32767)/float64(32768),0,-float64(32767)/float64(32768),0,float64(32767)/float64(32768),0,-float64(32767)/float64(32768),0,float64(32767)/float64(32768),0,-float64(32767)/float64(32768),0,float64(32767)/float64(32768),0}, f64) {
		t.Errorf("read badness %v", f64)
	}
	f.ReadItems(f32)
	if !floatEqual([]float32{-float32(32767)/float32(32768),0,float32(32767)/float32(32768),0,-float32(32767)/float32(32768),0,float32(32767)/float32(32768),0,-float32(32767)/float32(32768),0,float32(32767)/float32(32768),0,-float32(32767)/float32(32768),0,float32(32767)/float32(32768),0}, f32) {
		t.Errorf("read badness %v", f32)
	}
	ok := f.SetDoubleNormalization(false)
	if !ok {
		t.Error("expected previous norm mode to be true, was false")
	}
	f.SetFloatNormalization(false)
	ok = f.GetDoubleNormalization()
	if ok {
		t.Error("expected norm mode to be false, was true")
	}
	ok = f.GetFloatNormalization()
	if ok {
		t.Error("expected float norm mode to be false, was true")
	}
	f.ReadItems(f64)
	if !reflect.DeepEqual([]float64{-1,2,-3,4,-5,6,-7,8,-9,10,-11,12,-13,14,-15,16}, f64) {
		t.Errorf("read badness %v", f64)
	}
	f.ReadItems(f32)
	if !reflect.DeepEqual([]float32{-1,2,-3,4,-5,6,-7,8,-9,10,-11,12,-13,14,-15,16}, f32) {
		t.Errorf("read badness %v", f32)
	}
	
}

func floatEqual(f1, f2 interface{}) bool {
	// assume []float<size>
	switch t := f1.(type) {
	case []float64:
		f1v := f1.([]float64)
		f2v := f2.([]float64)
		if len(f1v) != len(f2v) {
			return false
		}
		for i, f1c := range f1v {
			if math.Fabs(f1c - f2v[i]) > math.SmallestNonzeroFloat64 {
				fmt.Println(f1, f2v, math.Fabs(f1c - f2v[i]), math.SmallestNonzeroFloat64)
				return false
			}
		}
	case []float32:
		f1v := f1.([]float32)
		f2v := f2.([]float32)
		if len(f1v) != len(f2v) {
			return false
		}
		for i, f1c := range f1v {
			if math.Fabs(float64(f1c - f2v[i])) > math.SmallestNonzeroFloat32 {
				return false
			}
		}
	default:
		return false
	}
	return true
}