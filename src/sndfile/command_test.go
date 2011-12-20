package sndfile

import "os"
import "encoding/binary"
import "math"
import "reflect"
import "testing"
import "strings"
import "unsafe"

func TestGetLibVersion(t *testing.T) {
	s, _ := GetLibVersion()
	t.Log(s)
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
	t.Log("TestGetLogInfo output: ", s)
	if err != nil {
		t.Error("TestGetLogInfo err: ", err)
	}
}

func TestFileCommands(t *testing.T) {
	var i Info
	f, err := Open("test/funky.aiff", Read, &i)
	if err != nil {
		t.Fatalf("open file failed %s", err)
	}

	max, err := f.CalcSignalMax()
	if err != nil {
		t.Fatalf("signal max failed %s", err)
	}

	t.Logf("max signal %f\n", max)

	max, err = f.CalcNormSignalMax()
	if err != nil {
		t.Fatalf("norm signal max failed %s", err)
	}

	t.Logf("norm max signal %f\n", max)

	maxarr, err := f.CalcMaxAllChannels()
	if err != nil {
		t.Fatalf("max all chans failed %s", err)
	}

	t.Logf("max all chans signal %v\n", maxarr)

	maxarr, err = f.CalcNormMaxAllChannels()
	if err != nil {
		t.Fatalf("max all chans failed %s", err)
	}

	t.Logf("norm max all chans signal %v\n", maxarr)

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
	t.Log("--- Supported simple formats")
	for i := 0; i < simpleformats; i++ {
		f, name, ext, ok := GetSimpleFormat(i)
		t.Logf("%08x %s %s\n", f, name, ext)
		if !ok {
			t.Error("error from GetSimpleFormat()")
		}
	}

	t.Log("--- Supported formats")
	// following is straight from examples in libsndfile distribution
	majorcount := GetMajorFormatCount()
	subcount := GetSubFormatCount()
	for m := 0; m < majorcount; m++ {
		f, name, ext, ok := GetMajorFormatInfo(m)
		if ok {
			t.Logf("--- MAJOR 0x%08x %v Extension: .%v\n", f, name, ext)
			af, aname, aext, ok := GetFormatInfo(f)
			if !ok || f != af || aname != name || aext != ext {
				t.Error(f, "!=", af, name, "!=", aname, ext, "!=", aext)
			}

			for s := 0; s < subcount; s++ {
				sf, sname, sok := GetSubFormatInfo(s)
				asf, aname, _, ok := GetFormatInfo(sf)
				if !ok || sf != asf || aname != sname {
					t.Error("sub", sf, "!=", asf, sname, "!=", aname)
				}
				var i Info
				i.Channels = 1
				i.Format = Format(f | sf)
				if sok && FormatCheck(i) {
					t.Logf("   0x%08x %v %v\n", f|sf, name, sname)
				}
			}
		} else {
			t.Logf("no format for number %v\n", m)
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
	i := &Info{0, 44100, 1, SF_FORMAT_WAV | SF_FORMAT_PCM_16, 0, 0}
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
	i.Format = SF_FORMAT_AIFF | SF_FORMAT_PCM_24
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
	i.Format = SF_FORMAT_AIFF | SF_FORMAT_PCM_16

	f, err := Open("addpeakchunk1.aiff", Write, &i)
	if err != nil {
		t.Fatalf("couldn't open file %v", err)
	}

	f.SetAddPeakChunk(false)
	_, err = f.WriteItems([]int16{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8})
	if err != nil {
		t.Error("write err:", err)
	}
	f.Close()

	f, err = Open("addpeakchunk1.aiff", Read, &i)
	if err != nil {
		t.Fatalf("couldn't open file %v", err)
	}

	// calc signals
	r, err := f.CalcSignalMax()
	if err != nil {
		t.Fatal("couldn't calculate signal max", err)
	}
	if r != 8.0 {
		t.Errorf("Signal max was %v\n", r)
	}

	r, err = f.CalcNormSignalMax()
	if err != nil {
		t.Fatal("couldn't calculate signal max", err)
	}
	if r != float64(8)/float64(0x8000) {
		t.Errorf("Signal max was %v not %v\n", r, float64(8)/float64(0x7fff))
	}

	ra, err := f.CalcMaxAllChannels()
	if err != nil {
		t.Fatal("couldn't calculate signal max", err)
	}
	if !reflect.DeepEqual(ra, []float64{5.0, 6.0, 7.0, 8.0}) {
		t.Errorf("Signal max was %v\n", ra)
	}

	ra, err = f.CalcNormMaxAllChannels()
	if err != nil {
		t.Fatal("couldn't calculate signal max", err)
	}
	if !reflect.DeepEqual(ra, []float64{5.0 / float64(0x8000), 6.0 / float64(0x8000), 7.0 / float64(0x8000), 8.0 / float64(0x8000)}) {
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

	i.Format = SF_FORMAT_AIFF | SF_FORMAT_DOUBLE
	// repeat for peak chunk, making sure peak chunk returns same value
	f, err = Open("addpeakchunk1.aiff", Write, &i)
	if err != nil {
		t.Fatalf("couldn't open file %v", err)
	}

	f.SetAddPeakChunk(true)
	_, err = f.WriteItems([]int16{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8})
	if err != nil {
		t.Error("write err:", err)
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
	if !reflect.DeepEqual(ra, []float64{5.0, 6.0, 7.0, 8.0}) {
		t.Errorf("Signal max was %v\n", ra)
	}

}

func TestNormalization(t *testing.T) {
	var i Info
	i.Format = SF_FORMAT_AIFF | SF_FORMAT_PCM_16
	i.Channels = 1
	i.Samplerate = 44100
	f, err := Open("norm.aiff", Write, &i)
	if err != nil {
		t.Fatal("couldn't open write file", err)
	}

	// set write normalization on
	f.SetDoubleNormalization(true)
	f.SetFloatNormalization(true)
	_, err = f.WriteItems([]float64{-1, 0, 1, 0, -1, 0, 1, 0, -1, 0, 1, 0, -1, 0, 1, 0})
	if err != nil {
		t.Fatal("couldn't write to file", err)
	}
	_, err = f.WriteItems([]float32{-1, 0, 1, 0, -1, 0, 1, 0, -1, 0, 1, 0, -1, 0, 1, 0})
	if err != nil {
		t.Fatal("couldn't write to file", err)
	}

	// set write normalization off
	f.SetDoubleNormalization(false)
	f.SetFloatNormalization(false)
	_, err = f.WriteItems([]float64{-1, 2, -3, 4, -5, 6, -7, 8, -9, 10, -11, 12, -13, 14, -15, 16})
	if err != nil {
		t.Fatal("couldn't write to file", err)
	}
	_, err = f.WriteItems([]float32{-1, 2, -3, 4, -5, 6, -7, 8, -9, 10, -11, 12, -13, 14, -15, 16})
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
	if !reflect.DeepEqual([]float64{-32767, 0, 32767, 0, -32767, 0, 32767, 0, -32767, 0, 32767, 0, -32767, 0, 32767, 0}, f64) {
		t.Errorf("read badness %v", f64)
	}
	f.ReadItems(f32)
	if !reflect.DeepEqual([]float32{-32767, 0, 32767, 0, -32767, 0, 32767, 0, -32767, 0, 32767, 0, -32767, 0, 32767, 0}, f32) {
		t.Errorf("read badness %v", f32)
	}
	f.ReadItems(f64)
	if !reflect.DeepEqual([]float64{-1, 2, -3, 4, -5, 6, -7, 8, -9, 10, -11, 12, -13, 14, -15, 16}, f64) {
		t.Errorf("read badness %v", f64)
	}
	f.ReadItems(f32)
	if !reflect.DeepEqual([]float32{-1, 2, -3, 4, -5, 6, -7, 8, -9, 10, -11, 12, -13, 14, -15, 16}, f32) {
		t.Errorf("read badness %v", f32)
	}

	f.Seek(0, Set)
	f.SetDoubleNormalization(true)
	f.SetFloatNormalization(true)

	n, err := f.ReadItems(f64)
	if err != nil || n != 16 {
		t.Fatal("bad read", err, n)
	}
	if !floatEqual([]float64{-float64(32767) / float64(32768), 0, float64(32767) / float64(32768), 0, -float64(32767) / float64(32768), 0, float64(32767) / float64(32768), 0, -float64(32767) / float64(32768), 0, float64(32767) / float64(32768), 0, -float64(32767) / float64(32768), 0, float64(32767) / float64(32768), 0}, f64) {
		t.Errorf("read badness %v", f64)
	}
	f.ReadItems(f32)
	if !floatEqual([]float32{-float32(32767) / float32(32768), 0, float32(32767) / float32(32768), 0, -float32(32767) / float32(32768), 0, float32(32767) / float32(32768), 0, -float32(32767) / float32(32768), 0, float32(32767) / float32(32768), 0, -float32(32767) / float32(32768), 0, float32(32767) / float32(32768), 0}, f32) {
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
	if !reflect.DeepEqual([]float64{-1, 2, -3, 4, -5, 6, -7, 8, -9, 10, -11, 12, -13, 14, -15, 16}, f64) {
		t.Errorf("read badness %v", f64)
	}
	f.ReadItems(f32)
	if !reflect.DeepEqual([]float32{-1, 2, -3, 4, -5, 6, -7, 8, -9, 10, -11, 12, -13, 14, -15, 16}, f32) {
		t.Errorf("read badness %v", f32)
	}

}

func floatEqual(f1, f2 interface{}) bool {
	// assume []float<size>
	switch f1.(type) {
	case []float64:
		f1v := f1.([]float64)
		f2v := f2.([]float64)
		if len(f1v) != len(f2v) {
			return false
		}
		for i, f1c := range f1v {
			if math.Abs(f1c-f2v[i]) > math.SmallestNonzeroFloat64 {
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
			if math.Abs(float64(f1c-f2v[i])) > math.SmallestNonzeroFloat32 {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func TestScaleFactor(t *testing.T) {
	var i Info
	i.Format = SF_FORMAT_AIFF | SF_FORMAT_FLOAT
	i.Channels = 1
	i.Samplerate = 8000
	os.Remove("scalefactor.aiff")
	f, err := Open("scalefactor.aiff", ReadWrite, &i)
	if err != nil {
		t.Fatal("couldn't open scale factor out", err)
	}
	f.SetIntFloatScaleWrite(false)
	out := []int16{2, 2, 4, 4, -2, -2, -4, -4}
	n, err := f.WriteItems(out)
	if n != int64(len(out)) || err != nil {
		t.Error("couldn't write items", err)
	}
	f.SetIntFloatScaleWrite(true)
	n, err = f.WriteItems(out)
	if n != int64(len(out)) || err != nil {
		t.Error("couldn't write items", err)
	}

	in := make([]int16, 2)
	f.Seek(0, Set)
	f.SetFloatIntScaleRead(false)
	n, err = f.ReadItems(in)
	if err != nil || n != 2 {
		t.Error("couldn't read items!", n, err)
	}

	if !reflect.DeepEqual(in, []int16{2, 2}) {
		t.Error("bad read 1", in)
	}
	f.SetFloatIntScaleRead(true)
	n, err = f.ReadItems(in)
	if !reflect.DeepEqual(in, []int16{16384, 16384}) {
		t.Error("bad read 2", in)
	}
	n, err = f.ReadItems(in)
	if !reflect.DeepEqual(in, []int16{32767, 32767}) {
		t.Error("bad read 3", in)
	}
	f.SetFloatIntScaleRead(false)
	n, err = f.ReadItems(in)
	if !reflect.DeepEqual(in, []int16{-2, -2}) {
		t.Error("bad read 4", in)
	}
}

func checkLength(t *testing.T) int32 {
	f, err := os.Open("update")
	if err != nil {
		t.Fatal("couldn't open file", err)
	}
	f.Seek(40, os.SEEK_SET)
	var length int32
	err = binary.Read(f, binary.LittleEndian, &length)
	if err != nil {
		t.Fatal("couldn't read from file", err)
	}
	f.Close()
	return length
}

func TestUpdateHeader(t *testing.T) {
	var i Info
	i.Format = SF_FORMAT_WAV | SF_FORMAT_PCM_16
	i.Channels = 1
	i.Samplerate = 8000
	os.Remove("update")
	f, err := Open("update", ReadWrite, &i)
	if err != nil {
		t.Fatal("couldn't open update.aiff")
	}
	b := f.SetUpdateHeaderAuto(true) // this appears to just return what you passed it
	if !b {
		t.Error("couldn't set SetUpdateHeaderAuto")
	}
	l := checkLength(t)
	if l != 0 {
		t.Error("length was non-zero before any writes?!?")
	}
	out := []int16{1, 4, 5, 2, 45, 12, 35, 2, 3, 56, 345, 64, 456, 7, 345, 62, 4567, 34, 67, 34, 56, 34, 56, 3, 456}
	var totsize int
	for i := 0; i <= 100; i++ {
		f.WriteItems(out)
		l = checkLength(t)
		if l != int32((i+1)*len(out)*2) { // WAV size is in bytes, not samples
			t.Error("header didn't update?", l, "!=", (i+1)*len(out)*2)
		}
		totsize += len(out) * 2
	}
	b = f.SetUpdateHeaderAuto(false)
	if b {
		t.Error("couldn't set SetUpdateHeaderAuto to false")
	}
	for i := 0; i <= 100; i++ {
		f.WriteItems(out)
		nl := checkLength(t)
		if l != nl { // WAV size is in bytes, not samples
			t.Error("header updated when auto = false", l, "!=", nl)
		}
		totsize += len(out) * 2
	}
	f.UpdateHeaderNow()
	nl := checkLength(t)
	if int(nl) != totsize {
		t.Error("bad size?", nl, "!=", totsize)
	}

	f.Close()
}

func TestBroadcast(t *testing.T) {
	var i Info
	i.Format = SF_FORMAT_WAV | SF_FORMAT_PCM_16
	i.Channels = 1
	i.Samplerate = 8000

	f, err := Open("broadcast", Write, &i)
	if err != nil {
		t.Fatal("couldn't open broacast file for write", err)
	}

	var bi BroadcastInfo
	bi.Description = "gosndfile test data"
	bi.Originator = "republic of nynex"
	bi.Originator_reference = "http://hydrogenproject.com"
	bi.Origination_date = "2011/09/27"
	bi.Origination_time = "17:49"
	bi.Time_reference_low = 123456
	bi.Time_reference_high = 7891011
	bi.Version = 1 // libsndfile always writes a 1
	bi.Umid = "ummm"
	bi.Coding_history = make([]int8, 257) // we don't set coding history, libsndfile does that
	bi.Coding_history[255] = 0x7f
	bi.Coding_history[256] = 0x11
	f.SetBroadcastInfo(&bi)
	f.Close()

	f, err = Open("broadcast", Read, &i)
	if err != nil {
		t.Fatal("couldn't open broadcast for read", err)
	}
	bi2, ok := f.GetBroadcastInfo()
	if !ok {
		t.Error("error retrieving broadcast info", err)
	}
	if bi.Description != bi2.Description {
		t.Error("desc doesn't match \"" + bi.Description + "\" \"" + bi2.Description + "\"")
	}
	if !reflect.DeepEqual(bi2.Coding_history, []int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 65, 61, 80, 67, 77, 44, 70, 61, 56, 48, 48, 48, 44, 87, 61, 49, 54, 44, 77, 61, 109, 111, 110, 111, 44, 84, 61, 108, 105, 98, 115, 110, 100, 102, 105, 108, 101, 45, 49, 46, 48, 46, 50, 53, 13, 10, 0, 0}) {
		t.Error("coding history mismatch")
	}
	bi.Coding_history = []int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 65, 61, 80, 67, 77, 44, 70, 61, 56, 48, 48, 48, 44, 87, 61, 49, 54, 44, 77, 61, 109, 111, 110, 111, 44, 84, 61, 108, 105, 98, 115, 110, 100, 102, 105, 108, 101, 45, 49, 46, 48, 46, 50, 53, 13, 10, 0, 0}
	if !reflect.DeepEqual(&bi, bi2) {
		t.Error("deepequal fails", &bi, bi2)
	}
}

func TestInstrument(t *testing.T) {
	var i Info
	i.Format = SF_FORMAT_AIFF | SF_FORMAT_PCM_24
	i.Channels = 2
	i.Samplerate = 44100
	f, err := Open("musicenztrumentz.aiff", Write, &i)
	if err != nil {
		t.Fatal("couldn't open file", err)
	}
	var inst Instrument
	inst.Gain = 1
	inst.Basenote = 64
	inst.Detune = -4
	inst.Velocity[0] = 14
	inst.Velocity[1] = 1
	inst.Key[0] = 0
	inst.Key[1] = 0
	inst.LoopCount = 1
	inst.Loops[0].Mode = Alternating
	inst.Loops[0].Start = 0x123
	inst.Loops[0].End = 0x321
	inst.Loops[0].Count = 0
	ok := f.SetInstrument(&inst)
	if !ok {
		t.Error("SetInstrument failed")
	}
	out := make([]int32, 44100*2)
	for i, _ := range out {
		out[i] = int32((i / 8) << 24)
	}
	n, err := f.WriteItems(out)
	if err != nil {
		t.Error("couldn't write to file", n, err)
	}
	f.Close()
	f, err = Open("musicenztrumentz.aiff", Read, &i)
	newinst := f.GetInstrument()
	if inst.Gain != newinst.Gain ||
		inst.Basenote != newinst.Basenote ||
		inst.Detune != newinst.Detune ||
		!reflect.DeepEqual(inst.Velocity, newinst.Velocity) ||
		!reflect.DeepEqual(inst.Key, newinst.Key) ||
		inst.LoopCount != newinst.LoopCount {
		t.Errorf("inst and newinst did not match\n%v\n%v", inst, newinst)
	}

}

// don't run this. looks like the actual command is a no-op!
func testRawOffset(t *testing.T) {
	var i Info
	i.Format = SF_FORMAT_RAW | SF_FORMAT_PCM_S8
	i.Samplerate = 8000
	i.Channels = 1

	f, err := Open("rawtest", Write, &i)
	if err != nil {
		t.Fatal("Writing file failed", err)
	}
	for i := int16(-256); i <= 255; i++ {
		f.WriteItems([]int16{i << 8})
	}
	/*	n, err := f.WriteFrames([]int16{0x100, 0x200,0x300,0x400,0x500,0x600,0x700,0x7f00})
		if err != nil || n != 4 {
			t.Error("Writing file failed", err)
		}*/
	var n int64
	f.Close()

	f, err = Open("rawtest", Read, &i)
	if err != nil {
		t.Fatal("reading file failed", err)
	}
	err = f.SetRawStartOffset(4)
	if err != nil {
		t.Error("set raw start failed!", err)
	}
	buf := make([]int16, 4)
	n, err = f.ReadFrames(buf)
	if err != nil || n != 4 {
		t.Fatal("reading file failed", n, err)
	}
	if !reflect.DeepEqual(buf, []int16{0x500, 0x600, 0x700, 0x800}) {
		t.Error("bad stuff", buf, []int16{0x500, 0x600, 0x700, 0x800})
	}
}

func TestClipping(t *testing.T) {
	var i Info
	i.Format = SF_FORMAT_AIFF | SF_FORMAT_PCM_16
	i.Samplerate = 8000
	i.Channels = 1
	f, err := Open("cliptest.aiff", Write, &i)
	if err != nil {
		t.Fatal("opening file for write failed", err)
	}
	f.SetClipping(true)
	n, err := f.WriteItems([]float64{0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0})
	if n != 8 || err != nil {
		t.Error("problem writing", n, err)
	}
	f.SetClipping(false)
	n, err = f.WriteItems([]float64{0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0})
	if n != 8 || err != nil {
		t.Error("problem writing", n, err)
	}
	f.Close()
	f, err = Open("cliptest.aiff", Read, &i)
	if err != nil {
		t.Fatal("opening file for read failed", err)
	}
	in := make([]int16, 16)
	n, err = f.ReadItems(in)
	if n != 16 || err != nil {
		t.Error("problem reading", n, err)
	}
	gold := []int16{0x4000, 0x7fff, 0x7fff, 0x7fff, 0x7fff, 0x7fff, 0x7fff, 0x7fff, 0x4000, 0x7fff, -16386, -2, 0x3ffe, 32765, -16388, -4}
	if !reflect.DeepEqual(in, gold) {
		t.Error("read doesn't match\n", gold, "\n", in)
	}
	f.Close()
}

func TestAmbisonic(t *testing.T) {
	var i Info
	i.Format = SF_FORMAT_WAVEX | SF_FORMAT_PCM_32
	i.Samplerate = 8000
	i.Channels = 1
	f, err := Open("ambisonictest.wav", Write, &i)
	if err != nil {
		t.Fatal("couldn't open file for write", err)
	}
	res := f.WavexSetAmbisonic(AmbisonicBFormat)
	if res != AmbisonicBFormat {
		t.Error("couldn't set ambisonic format", res, AmbisonicBFormat)
	}
	c := make([]int32, 31)
	for i := range c {
		c[i] = (1 << uint(i))
		if i%2 != 0 {
			c[i] *= -1
		}
	}
	f.WriteItems(c)
	f.Close()
	f, err = Open("ambisonictest.wav", Read, &i)
	if err != nil {
		t.Fatal("couldn't open file for reading", err)
	}
	if i.Format&SF_FORMAT_WAVEX == 0 {
		t.Errorf("Wrong format on read %x expected bit %x to be set\n", i.Format, SF_FORMAT_WAVEX)
	}
	res = f.WavexGetAmbisonic()
	if res != AmbisonicBFormat {
		t.Errorf("Wrong ambisonic answer %d, expected %d\n", res, AmbisonicBFormat)
	}
	f.Close()

}

// how do i make sure vbr quality is passed along correctly?

// i need to create a file with loop info. AIFF only?

// embedded file. buh?
