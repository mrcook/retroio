package tzx

import (
	"fmt"

	"github.com/mrcook/tzxit/tape"
)

// HardwareType
// ID: 33h (51d)
// This blocks contains information about the hardware that the programs on this tape use.
// Please include only machines and hardware for which you are 100% sure that it either runs
// (or doesn't run) on or with, or you know it uses (or doesn't use) the hardware or special
// features of that machine.
// If the tape runs only on the ZX81 (and TS1000, etc.) then it clearly won't work on any
// Spectrum or Spectrum variant, so there's no need to list this information.
// If you are not sure or you haven't tested a tape on some particular machine/hardware
// combination then do not include it in the list.
type HardwareType struct {
	TypeCount uint8          // N BYTE     Number of machines and hardware types for which info is supplied
	Machines  []HardwareInfo // HWINFO[N]  List of machines and hardware
}

// HardwareInfo (HWINFO)
// The list of hardware types and IDs is found at the end of this file.
type HardwareInfo struct {
	Type        uint8 // BYTE  Hardware type
	Id          uint8 // BYTE  Hardware ID
	Information uint8 // BYTE  Hardware information
}

func (h *HardwareType) Read(file *tape.File) {
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
	return "Hardware"
}

// ToString returns a human readable string of the block data
func (h HardwareType) ToString() string {
	str := fmt.Sprintf("> %s:\n", h.Name())
	for _, m := range h.Machines {
		str += fmt.Sprintf("- Type: %02X - %s\n", m.Type, hardwareReferenceTypes[m.Type])
		str += fmt.Sprintf("  ID:   %02X - %s\n", m.Id, hardwareReferenceIDs[m.Type][m.Id])
		str += fmt.Sprintf("  Info: %02X - %s\n", m.Information, hardwareInfoIDs[m.Information])
	}
	return str
}

// Information detailing the relationship between a piece of software and the hardware.
var hardwareInfoIDs = map[uint8]string{
	0x00: "The tape RUNS on this machine or with this hardware, but may or may not use the hardware or special features of the machine.",
	0x01: "The tape USES the hardware or special features of the machine, such as extra memory or a sound chip.",
	0x02: "The tape RUNS but it DOESN'T use the hardware or special features of the machine.",
	0x03: "The tape DOESN'T RUN on this machine or with this hardware.",
}

// This is the list of all hardware types and hardware identification ID's that are used
// in the 'Hardware info' block.
//
// By default you don't have to write any of the information if the game is made for the
// ZX Spectrum and complies with the following:
//   - runs on ZX Spectrum 48K
//   - runs on, but doesn't use any of the special hardware of ZX Spectrum 128K
//   - doesn't run on ZX Spectrum 16K
//
// If, for example, a game works on BOTH ZX 48K and 128K, and uses the hardware of the
// 128K Spectrum, then you would just include the 128K Spectrum in the list (because by
// default it has to work on 48K too).
//
// If the game is 128K ONLY then you would include two entries: The game works on AND
// uses the hardware of a 128K Spectrum AND the game DOESN'T work on a 48K Spectrum.
//
// If the game works on both 48K and 128K Spectrum, but it only uses the sound chip (AY)
// of the 128K Spectrum and none of its extra memory then you would only include the
// entry that says that the game uses the 'Classic AY hardware (compatible with 128k ZXs)'.

var hardwareReferenceTypes = map[uint8]string{
	0x00: "Computers",
	0x01: "External storage",
	0x02: "ROM/RAM type add-ons",
	0x03: "Sound devices",
	0x04: "Joysticks",
	0x05: "Mice",
	0x06: "Other controllers",
	0x07: "Serial ports",
	0x08: "Parallel ports",
	0x09: "Printers",
	0x0a: "Modems",
	0x0b: "Digitizers",
	0x0c: "Network adapters",
	0x0d: "Keyboards & keypads",
	0x0e: "AD/DA converters",
	0x0f: "EPROM programmers",
	0x10: "Graphics",
}

var hardwareReferenceIDs = map[uint8]map[uint8]string{
	0x00: { // Computers
		0x00: "ZX Spectrum 16k",
		0x01: "ZX Spectrum 48k, Plus",
		0x02: "ZX Spectrum 48k ISSUE 1",
		0x03: "ZX Spectrum 128k +(Sinclair)",
		0x04: "ZX Spectrum 128k +2 (grey case)",
		0x05: "ZX Spectrum 128k +2A, +3",
		0x06: "Timex Sinclair TC-2048",
		0x07: "Timex Sinclair TS-2068",
		0x08: "Pentagon 128",
		0x09: "Sam Coupe",
		0x0a: "Didaktik M",
		0x0b: "Didaktik Gama",
		0x0c: "ZX-80",
		0x0d: "ZX-81",
		0x0e: "ZX Spectrum 128k, Spanish version",
		0x0f: "ZX Spectrum, Arabic version",
		0x10: "Microdigital TK 90-X",
		0x11: "Microdigital TK 95",
		0x12: "Byte",
		0x13: "Elwro 800-3 ",
		0x14: "ZS Scorpion 256",
		0x15: "Amstrad CPC 464",
		0x16: "Amstrad CPC 664",
		0x17: "Amstrad CPC 6128",
		0x18: "Amstrad CPC 464+",
		0x19: "Amstrad CPC 6128+",
		0x1a: "Jupiter ACE",
		0x1b: "Enterprise",
		0x1c: "Commodore 64",
		0x1d: "Commodore 128",
		0x1e: "Inves Spectrum+",
		0x1f: "Profi",
		0x20: "GrandRomMax",
		0x21: "Kay 1024",
		0x22: "Ice Felix HC 91",
		0x23: "Ice Felix HC 2000",
		0x24: "Amaterske RADIO Mistrum",
		0x25: "Quorum 128",
		0x26: "MicroART ATM",
		0x27: "MicroART ATM Turbo 2",
		0x28: "Chrome",
		0x29: "ZX Badaloc",
		0x2a: "TS-1500",
		0x2b: "Lambda",
		0x2c: "TK-65",
		0x2d: "ZX-97",
	},
	0x01: { // External storage
		0x00: "ZX Microdrive",
		0x01: "Opus Discovery",
		0x02: "MGT Disciple",
		0x03: "MGT Plus-D",
		0x04: "Rotronics Wafadrive",
		0x05: "TR-DOS (BetaDisk)",
		0x06: "Byte Drive",
		0x07: "Watsford",
		0x08: "FIZ",
		0x09: "Radofin",
		0x0a: "Didaktik disk drives",
		0x0b: "BS-DOS (MB-02)",
		0x0c: "ZX Spectrum +3 disk drive",
		0x0d: "JLO (Oliger) disk interface",
		0x0e: "Timex FDD3000",
		0x0f: "Zebra disk drive",
		0x10: "Ramex Millenia",
		0x11: "Larken",
		0x12: "Kempston disk interface",
		0x13: "Sandy",
		0x14: "ZX Spectrum +3e hard disk",
		0x15: "ZXATASP",
		0x16: "DivIDE",
		0x17: "ZXCF",
	},
	0x02: { // ROM/RAM type add-ons
		0x00: "Sam Ram",
		0x01: "Multiface ONE",
		0x02: "Multiface 128k",
		0x03: "Multiface +3",
		0x04: "MultiPrint",
		0x05: "MB-02 ROM/RAM expansion",
		0x06: "SoftROM",
		0x07: "1k",
		0x08: "16k",
		0x09: "48k",
		0x0a: "Memory in 8-16k used",
	},
	0x03: { // Sound devices
		0x00: "Classic AY hardware (compatible with 128k ZXs)",
		0x01: "Fuller Box AY sound hardware",
		0x02: "Currah microSpeech",
		0x03: "SpecDrum",
		0x04: "AY ACB stereo (A+C=left, B+C=right); Melodik",
		0x05: "AY ABC stereo (A+B=left, B+C=right)",
		0x06: "RAM Music Machine",
		0x07: "Covox",
		0x08: "General Sound",
		0x09: "Intec Electronics Digital Interface B8001",
		0x0a: "Zon-X AY",
		0x0b: "QuickSilva AY",
		0x0c: "Jupiter ACE",
	},
	0x04: { // Joysticks
		0x00: "Kempston",
		0x01: "Cursor, Protek, AGF",
		0x02: "Sinclair 2 Left (12345)",
		0x03: "Sinclair 1 Right (67890)",
		0x04: "Fuller",
	},
	0x05: { // Mice
		0x00: "AMX mouse",
		0x01: "Kempston mouse",
	},
	0x06: { // Other controllers
		0x00: "Trickstick",
		0x01: "ZX Light Gun",
		0x02: "Zebra Graphics Tablet",
		0x03: "Defender Light Gun",
	},
	0x07: { // Serial ports
		0x00: "ZX Interface 1",
		0x01: "ZX Spectrum 128k",
	},
	0x08: { // Parallel ports
		0x00: "Kempston S",
		0x01: "Kempston E",
		0x02: "ZX Spectrum +3",
		0x03: "Tasman",
		0x04: "DK'Tronics",
		0x05: "Hilderbay",
		0x06: "INES Printerface",
		0x07: "ZX LPrint Interface 3",
		0x08: "MultiPrint",
		0x09: "Opus Discovery",
		0x0a: "Standard 8255 chip with ports 31,63,95",
	},
	0x09: { // Printers
		0x00: "ZX Printer, Alphacom 32 & compatibles",
		0x01: "Generic printer",
		0x02: "EPSON compatible",
	},
	0x0a: { // Modems
		0x00: "Prism VTX 5000",
		0x01: "T/S 2050 or Westridge 2050",
	},
	0x0b: { // Digitizers
		0x00: "RD Digital Tracer",
		0x01: "DK'Tronics Light Pen",
		0x02: "British MicroGraph Pad",
		0x03: "Romantic Robot Videoface",
	},
	0x0c: { // Network adapters
		0x00: "ZX Interface 1",
	},
	0x0d: { // Keyboards & keypads
		0x00: "Keypad for ZX Spectrum 128k",
	},
	0x0e: { // AD/DA converters
		0x00: "Harley Systems ADC 8.2",
		0x01: "Blackboard Electronics",
	},
	0x0f: { // EPROM programmers
		0x00: "Orme Electronics",
	},
	0x10: { // Graphics
		0x00: "WRX Hi-Res",
		0x01: "G007",
		0x02: "Memotech",
		0x03: "Lambda Colour",
	},
}
