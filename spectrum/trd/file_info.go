package trd

import (
	"fmt"
	"github.com/mrcook/retroio/storage"
)

// FileInformation represents a single file information obtained from the disk catalog
type FileInformation struct {
	Filename        [8]byte
	FileType        FileType
	StartAddress    uint16
	LengthInBytes   uint16
	LengthInSectors uint8
	StartingSector  uint8
	StartingTrack   uint8
}

func (i *FileInformation) Read(reader *storage.Reader) error {
	if _, err := reader.Peek(16); err != nil {
		return err
	}

	copy(i.Filename[:], reader.ReadBytes(8))

	ext := reader.ReadByte()
	switch ext {
	case 'b':
		fallthrough
	case 'B':
		i.FileType = &FileTypeBasic{}
	case 'c':
		fallthrough
	case 'C':
		i.FileType = &FileTypeCode{}
	default:
		i.FileType = &FileTypeOther{Extension: ext}
	}

	i.StartAddress = reader.ReadShort()
	i.LengthInBytes = reader.ReadShort()
	i.LengthInSectors = reader.ReadByte()
	i.StartingSector = reader.ReadByte()
	i.StartingTrack = reader.ReadByte()

	return nil
}

// IsDeleted returns whether the given entry represents a deleted file
// The procedure for deleting a file is to replace the first byte of its name with the code 0x00 or 0x01.
func (i FileInformation) IsDeleted() bool {
	switch i.Filename[0] {
	case 0x00:
		fallthrough
	case 0x01:
		return true
	}

	return false
}

func (i FileInformation) String() string {
	str := fmt.Sprintf("%s\n", i.FileType.Name())
	str += fmt.Sprintf(" - Filename:    %s\n", i.Filename)
	str += fmt.Sprintf(" - Sectors:     %d\n", i.LengthInSectors)
	str += fmt.Sprintf(" - Bytes:       %d\n", i.LengthInBytes)
	str += i.FileType.Info(i)

	return str
}
