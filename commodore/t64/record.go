package t64

import "fmt"

// File record
// - C64s filetype
//        0 = free entry
//        1 = normal tape file
//        2 = tape file with header: header is saved just before file data
//        3 = memory snapshot v0.9, uncompressed
//        4 = tape block
//        5 = digitized stream
//        6 ... 255 reserved
//   VICE only supports file type 1.
//   Types 3, 4 and 5 are subject to change (and are rarely used).
// - 1541 file type (0x82 for PRG, 0x81 for SEQ, etc). You will find it can vary between 0x01, 0x44, and the normal
//   D64 values. In reality any value that is not a $00 is seen as a PRG file. When this value is a $00 (and the
//   previous byte at "filetype" is >1), then the file is a special T64 "FRZ" (frozen) C64s session snapshot.
// - Start address (or Load address). This is the first two bytes of the C64 file which is usually the load address
//   (typically $01 $08). If the file is a snapshot,the address will be 0.
// - End address (actual end address in memory, if the file was loaded into a C64). If the file is a snapshot, then
//   the address will be a 0.
// - Not used
// - Offset into the container file (from the beginning) of where the C64 file starts
// - Not used
// - C64 filename (in PETASCII, padded with $20, not $A0)
type Record struct {
	Type         uint8    // C64s entry type
	FileType     uint8    // 1541 file type: PRG: 0x82, SEQ: 0x81, etc. (USR ?)
	StartAddress uint16   // Start address (or Load address)
	EndAddress   uint16   // End address (actual end address in memory if the file was loaded into a C64).
	Unused1      uint16   // unused value
	Offset       uint32   // Offset of file contents start within T64 file
	Unused2      uint32   // unused value
	Filename     [16]byte // C64 filename, in PETASCII, padded with $20
}

func (r Record) String() string {
	str := ""
	str += fmt.Sprintf("Filename:      %s\n", r.Filename)
	str += fmt.Sprintf("Type:          %s: %s\n", r.fileTypeLabel(r.FileType), r.entryTypeLabel(r.Type))
	str += fmt.Sprintf("Start Address: %d\n", r.StartAddress)
	str += fmt.Sprintf("End Address:   %d\n", r.EndAddress)
	str += fmt.Sprintf("Data length:   %d\n", r.EndAddress-r.StartAddress)
	str += fmt.Sprintf("T64 Offset:    %d\n", r.Offset)
	return str
}

func (r Record) entryTypeLabel(id byte) string {
	var label string
	switch id {
	case 0x00:
		label = "Free Entry"
	case 0x01:
		label = "Normal tape file"
	case 0x02:
		label = "Tape file with header"
	case 0x03:
		label = "Memory Snapshot v0.9 (uncompressed)"
	case 0x04:
		label = "Tape Block"
	case 0x05:
		label = "Digitized Stream"
	default:
		label = "Unknown reserved type"
	}
	return label
}

func (r Record) fileTypeLabel(id byte) string {
	var label string
	switch id {
	case 0x81:
		label = "SEQ"
	case 0x82:
		label = "PRG"
	default:
		label = "???" // what about USR value?
	}
	return label
}
