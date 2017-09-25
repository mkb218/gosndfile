// +build !legacy

package sndfile

// #cgo pkg-config: sndfile
// #include <stdlib.h>
// #include <sndfile.h>
// #include <string.h>
import "C"
import "unsafe"

func broadcastFromC(c *C.SF_BROADCAST_INFO) *BroadcastInfo {
	bi := new(BroadcastInfo)
	bi.Description = trim(C.GoStringN(&c.description[0], C.int(len(c.description[:]))))
	bi.Originator = trim(C.GoStringN(&c.originator[0], C.int(len(c.originator[:]))))
	bi.Originator_reference = trim(C.GoStringN(&c.originator_reference[0], C.int(len(c.originator_reference[:]))))
	bi.Origination_date = trim(C.GoStringN(&c.origination_date[0], C.int(len(c.origination_date[:]))))
	bi.Origination_time = trim(C.GoStringN(&c.origination_time[0], C.int(len(c.origination_time[:]))))
	bi.Time_reference_low = uint32(c.time_reference_low)
	bi.Time_reference_high = uint32(c.time_reference_high)
	bi.Version = uint16(c.version)
	bi.Umid = trim(C.GoStringN(&c.umid[0], C.int(len(c.umid[:]))))
	coding_history_bytes := make([]byte, 0, c.coding_history_size)
	for i, r := range c.coding_history {
		if i >= int(c.coding_history_size) {
			break
		}
		coding_history_bytes = append(coding_history_bytes, byte(r))
	}
	return bi
}

func cFromBroadcast(bi *BroadcastInfo) (c *C.SF_BROADCAST_INFO) {
	c = new(C.SF_BROADCAST_INFO)
	arrFromGoString(c.description[:], bi.Description)
	arrFromGoString(c.originator[:], bi.Originator)
	arrFromGoString(c.originator_reference[:], bi.Originator_reference)
	arrFromGoString(c.origination_date[:], bi.Origination_date)
	arrFromGoString(c.origination_time[:], bi.Origination_time)
	c.time_reference_low = C.uint32_t(bi.Time_reference_low)
	c.time_reference_high = C.uint32_t(bi.Time_reference_high)
	c.version = C.short(bi.Version)
	arrFromGoString(c.umid[:], bi.Umid)
	ch := bi.Coding_history
	if len(bi.Coding_history) > 256 {
		ch = bi.Coding_history[0:256]
	}
	c.coding_history_size = C.uint32_t(len(ch))
	for i, r := range ch {
		c.coding_history[i] = C.char(r)
	}
	return c
}

// Set instrument information from file including MIDI base note, keyboard mapping and looping information (start/stop and mode).

// Return true if the file header contains instrument information for the file. false otherwise.
func (f *File) SetInstrument(i *Instrument) bool {
	c := new(C.SF_INSTRUMENT)
	c.gain = C.int(i.Gain)
	c.basenote = C.char(i.Basenote)
	c.detune = C.char(i.Detune)
	c.velocity_lo = C.char(i.Velocity[0])
	c.velocity_hi = C.char(i.Velocity[1])
	c.key_lo = C.char(i.Key[0])
	c.key_hi = C.char(i.Key[1])
	c.loop_count = C.int(i.LoopCount)
	var index int
	for ; index < i.LoopCount; index++ {
		c.loops[index].mode = C.int(i.Loops[index].Mode)
		c.loops[index].start = C.uint32_t(i.Loops[index].Start)
		c.loops[index].end = C.uint32_t(i.Loops[index].End)
		c.loops[index].count = C.uint32_t(i.Loops[index].Count)
	}
	for ; index < 16; index++ {
		c.loops[index].mode = C.int(None)
		// why is this necessary? libsndfile doesn't check loopcount for AIFF
	}

	r := C.sf_command(f.s, C.SFC_SET_INSTRUMENT, unsafe.Pointer(c), C.int(unsafe.Sizeof(*c)))
	return (r == C.SF_TRUE)
}