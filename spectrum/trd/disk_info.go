package trd

import (
	"fmt"
	"github.com/mrcook/retroio/storage"
)

// DiskInformation stores the disk information obtained from the 9th sector of the system track
type DiskInformation struct {
	NextFreeSector     uint8
	NextFreeTrack      uint8
	DiskType           DiskType
	NumFiles           uint8
	NumFreeSectors     uint16
	NumDeletedFiles    uint8
	Label              [8]byte
}

func (di *DiskInformation) Read(reader *storage.Reader) error {
	if _, err := reader.Peek(256); err != nil {
		return err
	}

	// Skip to position 225
	reader.ReadBytes(225)

	di.NextFreeSector = reader.ReadByte()
	di.NextFreeTrack = reader.ReadByte()

	di.DiskType = DiskType(reader.ReadByte())

	di.NumFiles = reader.ReadByte()
	di.NumFreeSectors = reader.ReadShort()

	// Skip the number of sectors per track (MAIN_BYTE) which is always the same
	reader.ReadByte()

	// Skip two zero bytes
	reader.ReadShort()

	// Skip nine whitespaces
	reader.ReadBytes(9)

	// Skip one more zero byte
	reader.ReadByte()

	di.NumDeletedFiles = reader.ReadByte()
	copy(di.Label[:], reader.ReadBytes(8))

	return nil
}

func (di DiskInformation) String() string {
	str := ""
	str += fmt.Sprintf("Type:          %s\n", di.DiskType)
	str += fmt.Sprintf("Label:         %s\n", di.Label)
	str += fmt.Sprintf("Total files:   %d\n", di.NumFiles)
	str += fmt.Sprintf("Deleted files: %d\n", di.NumDeletedFiles)
	str += fmt.Sprintf("Free sectors:  %d\n", di.NumFreeSectors)

	return str
}
