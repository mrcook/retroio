package disk

// Disk DOS version codes.
// DOS version $41 is the most common, being used for the 1540/41/71 drives,
// and other variants such as SpeedDOS and Dolphin-DOS.
// Any other value enables "soft write protection" (see `docs.md`)
var DosVersions = map[uint8]string{
	0x00: "Not Set",
	0x41: "('A') CBM DOS v2.6 (1540/41/71), and others",
	0x44: "('D') CBM DOS v3.x (1581)",
	0x50: "('P') PrologicDOS 1541 and ProSpeed 1571 2.0",
}

// DOS type codes.
// NOTE: the 1581 is possibly CBM DOS v10.0.
var DosTypes = map[string]string{
	"2A": "CBM DOS v2.6",
	"2P": "PrologicDOS 1541 and ProSpeed 1571 2.0",
	"3D": "CBM DOS v3.x",
	"4A": "Professional DOS Release",
}
