// The sector information header follows closely the NEC µPD765 specification.
// https://web.archive.org/web/20170320235834/http://dec8.info/Apple/uPD72070.pdf
package dsk

import (
	"encoding/binary"
	"fmt"

	"retroio/storage"
)

const sectorInformationBlockSize = 8 // bytes

// Sector information block
//
// * For 8k Sectors (N="6"), only 1800h bytes is stored.
type SectorInformation struct {
	Track  uint8  // C   Cylinder Number is the current/selected track number: 0 through 76.
	Side   uint8  // H   Head Address is the head number: 0 or 1
	ID     uint8  // R   Record / sector number
	Size   uint8  // N   Number of data bytes written to sector
	ST1    uint8  // ST1 Error Status Register 1
	ST2    uint8  // ST2 Error Status Register 2
	Unused uint16 // not used (0)
}

// Read the track information header.
func (s *SectorInformation) Read(reader *storage.Reader) error {
	return binary.Read(reader, binary.LittleEndian, s)
}

// DataRead reads the data from the disk
func (s *SectorInformation) DataRead(reader *storage.Reader) ([]byte, error) {
	if s.Size > 3 {
		return nil, fmt.Errorf("unknown sector size value 0x%02X", s.Size)
	}

	data := make([]byte, s.SectorByteSize())
	err := binary.Read(reader, binary.LittleEndian, data)

	return data, err
}

// SectorByteSize calculates how many bytes the sector size is.
func (s SectorInformation) SectorByteSize() int {
	bytes := 0
	switch s.Size {
	case 0:
		bytes = 128
	case 1:
		bytes = 256
	case 2:
		bytes = 512
	case 3:
		bytes = 1024
	default:
		return -1
	}
	return bytes
}

func (s SectorInformation) String() string {
	str := ""
	str += fmt.Sprintf("ID:    %d\n", s.ID)
	str += fmt.Sprintf("Side:  %d\n", s.Side)
	str += fmt.Sprintf("Track: %d\n", s.Track)
	str += fmt.Sprintf("Size:  %d (%d bytes)\n", s.Size, s.SectorByteSize())
	str += fmt.Sprintf("ST1:   0x%02X: %s\n", s.ST1, s.st1Label())
	str += fmt.Sprintf("ST2:   0x%02X: %s\n", s.ST2, s.st2Label())
	return str
}

// NEC765 specification for "Status Register 1"
//
// EN: When the FDC tries to access a Sector beyond  the final Sector of a Cylinder, this flag is set.
// DE: When the FDC detects a CRC error in either the ID field or the data field, this flag is set.
// OR: If the FDC is not serviced by the main-systems during data transfers, with in a certain time interval, this flag is set.
// ND: During execution of READ DATA, WRITE DELETED DATA or SCAN Command, if the FDC cannot find the Sector specified in the IDR Register, this flag is set.
//     During executing the READ ID Command, if the FDC cannot read the ID field without an error, then this flag is set.
//     During the execution of the READ A Cylinder Command, if the starting sec
// NW: During execution of WRITE DATA, WRITE DELETED DATA or Format A Cylinder Command, if the FDC detects a write protect signal from the FDD, then this flag is set.
// MA: If the FDC cannot detect the ID Address Mark after encountering the index hole twice, then th is flag is set.
//     If the FDC cannot detect the Data Address Mark or Deleted Data Address Mark, this flag is set. Also at the same time, the MD (Missing Address Mark in Data Field) of Status Register 2 is set.
func (s SectorInformation) st1Label() string {
	var label string

	switch s.ST1 {
	case 0xB7:
		label = "EN (End of Cylinder)"
	case 0xB5:
		label = "DE (Data Error)"
	case 0xB4:
		label = "OR (Over Run)"
	case 0xB2:
		label = "ND (No Data)"
	case 0xB1:
		label = "NW (Not Writable)"
	case 0xB0:
		label = "MA (Missing Address Mark)"
	default:
		label = "unknown"
	}

	return label
}

// NEC765 specification for "Status Register 2"
//
// CM: During executing the READ DATA or SCAN Command, if the FDC encounters a Sector which contains a Deleted Data Address Mark, this flag is set.
// DD: If the FDC detects a CRC error in the data field then this flag is set.
// WC: This bit is related with the ND bit, and when the contents of C on the medium is different from that stored in the IDR, this flag is set.
// SH: During execution, the SCAN Command, if the condition of "equal" is satisfied, this flag is set.
// SN: During executing the SCAN Command, if the FDC cannot find a Sector on the cylinder which meets the condition, then this flag is set.
// BC: This bit is related with the ND bit, and when the content of C on the medium is different from that stored in the IDR and the content of C is FF, then this flag is set.
// MD: When data is read from the medium, if the FDC cannot find a Data Address Mark or Deleted Data Address Mark, then this flag is set.
func (s SectorInformation) st2Label() string {
	var label string

	switch s.ST2 {
	case 0xB6:
		label = "CM (Control Mark)"
	case 0xB5:
		label = "DD (Data Error in Data field)"
	case 0xB4:
		label = "WC (Wrong Cylinder)"
	case 0xB3:
		label = "SH (Scan Equal Hit)"
	case 0xB2:
		label = "SN (Scan Not Satisfied)"
	case 0xB1:
		label = "BC (Bad Cylinder)"
	case 0xB0:
		label = "MD (Missing address Mark in Data field)"
	default:
		label = "unknown"
	}

	return label
}
