package t64

import (
	"encoding/binary"
	"io"

	"github.com/mrcook/retroio/librio/storage"
)

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

// Read the record header data
func (r *Record) Read(reader *storage.Reader) error {
	return binary.Read(reader, binary.LittleEndian, r)
}

// DataRead reads the data for the record
func (r *Record) DataRead(reader *storage.Reader, dataOffset int) ([]byte, error) {
	length := int(r.EndAddress - r.StartAddress)

	data := make([]byte, length)
	_, err := reader.Read(data)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return data, nil
}
