// +build legacy

package sndfile

// #cgo pkg-config: sndfile
// #include <stdlib.h>
// #include <sndfile.h>
// #include <string.h>
import "C"

func broadcastFromC(c *C.SF_BROADCAST_INFO) *BroadcastInfo {
	bi := new(BroadcastInfo)
	bi.Description = trim(C.GoStringN(&c.description[0], C.int(len(c.description[:]))))
	bi.Originator = trim(C.GoStringN(&c.originator[0], C.int(len(c.originator[:]))))
	bi.Originator_reference = trim(C.GoStringN(&c.originator_reference[0], C.int(len(c.originator_reference[:]))))
	bi.Origination_date = trim(C.GoStringN(&c.origination_date[0], C.int(len(c.origination_date[:]))))
	bi.Origination_time = trim(C.GoStringN(&c.origination_time[0], C.int(len(c.origination_time[:]))))
	bi.Time_reference_low = uint32(uint(c.time_reference_low))
	bi.Time_reference_high = uint32(uint(c.time_reference_high))
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
