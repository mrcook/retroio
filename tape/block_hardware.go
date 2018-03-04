package tape

import "fmt"

// HardwareType
// ID: 33h (51d)
// This blocks contains information about the hardware that the programs on this tape use.
// Please include only machines and hardware for which you are 100% sure that it either runs
// (or doesn't run) on or with, or you know it uses (or doesn't use) the hardware or special
// features of that machine.
type HardwareType struct {
	TypeCount uint8          // N BYTE     Number of machines and hardware types for which info is supplied
	Machines  []HardwareInfo // HWINFO[N]  List of machines and hardware
}

// HardwareInfo
// A list of hardware types and ID and sub-types can be found in the TZX specification
// (https://www.worldofspectrum.org/TZXformat.html)
type HardwareInfo struct {
	Type        uint8 // BYTE  Hardware type
	Id          uint8 // BYTE  Hardware ID
	Information uint8 // BYTE  Hardware information:
	//                           00 - The tape RUNS on this machine or with this hardware,
	//                                but may or may not use the hardware or special features of the machine.
	//                           01 - The tape USES the hardware or special features of the machine,
	//                                such as extra memory or a sound chip.
	//                           02 - The tape RUNS but it DOESN'T use the hardware
	//                                or special features of the machine.
	//                           03 - The tape DOESN'T RUN on this machine or with this hardware.
}

func (h *HardwareType) Process(file *File) {
	h.TypeCount, _ = file.ReadByte()

	for i := 0; i < int(h.TypeCount); i++ {
		var m HardwareInfo
		m.Type, _ = file.ReadByte()
		m.Id, _ = file.ReadByte()
		m.Information, _ = file.ReadByte()
		h.Machines = append(h.Machines, m)
	}
}

func (h HardwareType) Id() int {
	return 51
}

func (h HardwareType) Name() string {
	return "Hardware Type and Information"
}

// Metadata returns a human readable string of the block data
func (h HardwareType) Metadata() string {
	str := ""
	str += fmt.Sprintf("Display Time: %d\n", h.TypeCount)

	str += fmt.Sprintf("Machines:\n")
	for _, b := range h.Machines {
		str += fmt.Sprintf("- Type: %d\n", b.Type)
		str += fmt.Sprintf("  ID:   %d\n", b.Id)
		str += fmt.Sprintf("  Info: %d\n", b.Information)
	}
	return str
}
