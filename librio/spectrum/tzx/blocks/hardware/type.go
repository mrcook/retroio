package hardware

// Type of Hardware being represented
type HwType uint8

const (
	Computer        HwType = 0x00 // "Computers"
	ExternalStorage HwType = 0x01 // "External storage"
	RomRamAddon     HwType = 0x02 // "ROM/RAM type add-ons"
	SoundDevice     HwType = 0x03 // "Sound devices"
	Joystick        HwType = 0x04 // "Joysticks"
	Mice            HwType = 0x05 // "Mice"
	OtherController HwType = 0x06 // "Other controllers"
	SerialPort      HwType = 0x07 // "Serial ports"
	ParallelPort    HwType = 0x08 // "Parallel ports"
	Printer         HwType = 0x09 // "Printers"
	Modem           HwType = 0x0a // "Modems"
	Digitizer       HwType = 0x0b // "Digitizers"
	NetworkAdapter  HwType = 0x0c // "Network adapters"
	Keyboard        HwType = 0x0d // "Keyboards & keypads"
	AddaConverter   HwType = 0x0e // "AD/DA converters"
	EpromProgrammer HwType = 0x0f // "EPROM programmers"
	Graphics        HwType = 0x10 // "Graphics"
)
