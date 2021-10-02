// The sector information header follows closely the NEC µPD765 specification.
// https://web.archive.org/web/20170320235834/http://dec8.info/Apple/uPD72070.pdf
package dsk

import (
	"encoding/binary"
	"fmt"

	"github.com/mrcook/retroio/librio/storage"
)

// Sector information block
//
// Closely follows the NEC765 specification.
//
// * N = 0x06: 8k Sectors, only 1800h bytes are stored.
// * N = 0x07: 16K will be stored
// * N = 0x08: is equivalent to N=0x00 as only 3 bits are used (see docs.md)
type SectorInformation struct {
	Track uint8 // C   Cylinder Number is the current/selected track number: 0 through 76.
	Side  uint8 // H   Head Address is the head number: 0 or 1
	ID    uint8 // R   Record / sector number
	Size  uint8 // N   Number of data bytes written to sector (enum 0-7)
	ST1   uint8 // ST1 Error Status Register 1
	ST2   uint8 // ST2 Error Status Register 2

	// Used only with EXTENDED disk images.
	//
	// The location of each sectors data is found by adding the size of the
	// previous sectors, plus the size of the 256 byte header.
	SectorDataLength uint16 // little endian notation.
}

// Read the track information header.
func (s *SectorInformation) Read(reader *storage.Reader) error {
	return binary.Read(reader, binary.LittleEndian, s)
}

// dataRead reads the data from the disk
func (s *SectorInformation) dataRead(reader *storage.Reader) ([]byte, error) {
	if s.Size > 3 {
		return nil, fmt.Errorf("unknown sector size value 0x%02X", s.Size)
	}

	sectorSize, ok := SectorSizeMap[s.Size]
	if !ok {
		return nil, fmt.Errorf("invalid sector size byte")
	}

	data := make([]byte, sectorSize)
	err := binary.Read(reader, binary.LittleEndian, data)

	return data, err
}
