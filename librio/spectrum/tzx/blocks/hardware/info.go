package hardware

// Info (HWINFO)
// The list of hardware types and IDs is found at the end of this file.
type Info struct {
	Type          HwType        // Hardware type
	Id            uint8         // Hardware ID
	Compatibility Compatibility // Hardware compatibility flag
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

var HardwareReferenceTypes = map[HwType]string{
	Computer:        "Computers",
	ExternalStorage: "External storage",
	RomRamAddon:     "ROM/RAM type add-ons",
	SoundDevice:     "Sound devices",
	Joystick:        "Joysticks",
	Mice:            "Mice",
	OtherController: "Other controllers",
	SerialPort:      "Serial ports",
	ParallelPort:    "Parallel ports",
	Printer:         "Printers",
	Modem:           "Modems",
	Digitizer:       "Digitizers",
	NetworkAdapter:  "Network adapters",
	Keyboard:        "Keyboards & keypads",
	AddaConverter:   "AD/DA converters",
	EpromProgrammer: "EPROM programmers",
	Graphics:        "Graphics",
}

var HardwareReferenceIDs = map[HwType]map[uint8]string{
	Computer: {
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
	ExternalStorage: {
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
	RomRamAddon: {
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
	SoundDevice: {
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
	Joystick: {
		0x00: "Kempston",
		0x01: "Cursor, Protek, AGF",
		0x02: "Sinclair 2 Left (12345)",
		0x03: "Sinclair 1 Right (67890)",
		0x04: "Fuller",
	},
	Mice: {
		0x00: "AMX mouse",
		0x01: "Kempston mouse",
	},
	OtherController: {
		0x00: "Trickstick",
		0x01: "ZX Light Gun",
		0x02: "Zebra Graphics Tablet",
		0x03: "Defender Light Gun",
	},
	SerialPort: {
		0x00: "ZX Interface 1",
		0x01: "ZX Spectrum 128k",
	},
	ParallelPort: {
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
	Printer: {
		0x00: "ZX Printer, Alphacom 32 & compatibles",
		0x01: "Generic printer",
		0x02: "EPSON compatible",
	},
	Modem: {
		0x00: "Prism VTX 5000",
		0x01: "T/S 2050 or Westridge 2050",
	},
	Digitizer: {
		0x00: "RD Digital Tracer",
		0x01: "DK'Tronics Light Pen",
		0x02: "British MicroGraph Pad",
		0x03: "Romantic Robot Videoface",
	},
	NetworkAdapter: {
		0x00: "ZX Interface 1",
	},
	Keyboard: {
		0x00: "Keypad for ZX Spectrum 128k",
	},
	AddaConverter: {
		0x00: "Harley Systems ADC 8.2",
		0x01: "Blackboard Electronics",
	},
	EpromProgrammer: {
		0x00: "Orme Electronics",
	},
	Graphics: {
		0x00: "WRX Hi-Res",
		0x01: "G007",
		0x02: "Memotech",
		0x03: "Lambda Colour",
	},
}
