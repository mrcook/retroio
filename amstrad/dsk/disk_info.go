package dsk

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/mrcook/retroio/storage"
)

// Disc Information Block
//
// Always at offset 0 in the disk image file.
// All tracks must be the same size.
//
// "MV - CPC" or "EXTENDED CPC" must be present, because it is used to identify the file as a disk image.
//   - "\r" is the C programming language equivalent of ASCII character 13.
//   - "\n" is the C programming language equivalent of ASCII character 10.
//
// "size of track" includes the 256 byte (&0100) Track Information Block.
type DiskInformation struct {
	// "MV - CPCEMU Disk-File\r\nDisk-Info\r\n" for standard Amstrad disks.
	// "EXTENDED CPC DSK File\r\nDisk-Info\r\n" for EXTENDED DiSK Format.
	Identifier [34]byte

	Creator [14]byte // name of the creator (utility, or emulator)
	Tracks  uint8    // number of tracks
	Sides   uint8    // number of sides

	// Size of a track (STANDARD disks) is used to calculate the location of the data for a chosen track.
	// Unused for EXTENDED disk images
	TrackSize uint16

	// The following bytes are only used with EXTENDED disks.
	// For STANDARD disk images this is unused padding.
	//
	// * Depending on the information in the disk information block, the table
	//   contains track lengths for a single/doubled sided floppy discs.
	// * Track lengths are stored in the same order as the tracks in the disk image.
	// * A size of "0" indicates an unformatted track - no data, no track info block.
	// * Length of track data = (high byte of track length) * 256
	// * Length includes the size of the TRACK INFORMATION BLOCK (256 bytes)
	//
	// NOTE:
	// The location of a Track Information Block for a chosen track is found by
	// summing the sizes of all tracks up to the chosen track plus the size of
	// the Disc Information Block (&100 bytes).
	// in the disc image.
	TrackSizeTable [204]uint8 // to Track Info block at 0x0100
}

// Read the disk information header.
func (d *DiskInformation) Read(reader *storage.Reader) error {
	return binary.Read(reader, binary.LittleEndian, d)
}

// "MV - CPCEMU Disk-File\r\nDisk-Info\r\n" for standard Amstrad disks.
func (d DiskInformation) isStandardDisk() bool {
	return strings.HasPrefix(reformatIdentifier(d.Identifier[:]), "MV - CPCEMU Disk-File")
}

// "EXTENDED CPC DSK File\r\nDisk-Info\r\n" for EXTENDED DiSK images.
func (d DiskInformation) isExtendedDisk() bool {
	return strings.HasPrefix(reformatIdentifier(d.Identifier[:]), "EXTENDED CPC DSK File")
}

// Amstrad disc media type (sidedness)
// See `docs.md` for more information on the type value.
func (d *DiskInformation) mediaType() uint8 {
	if d.Sides == 2 {
		// double sided
		return 1
	}

	// single sided
	return 0
}

func (d DiskInformation) String() string {
	str := ""
	str += fmt.Sprintf("Identifier: %s\n", reformatIdentifier(d.Identifier[:]))
	str += fmt.Sprintf("Creator:    %s\n", d.printableCreatorLabel())
	str += fmt.Sprintf("Tracks:     %d\n", d.Tracks)
	str += fmt.Sprintf("Sides:      %d\n", d.Sides)

	if d.isStandardDisk() {
		str += fmt.Sprintf("Track Size: %d\n", d.TrackSize)
	}

	return str
}

func (d DiskInformation) printableCreatorLabel() []byte {
	var i int
	for _, c := range d.Creator {
		if c == 0 {
			break
		}
		i++
	}
	return d.Creator[:i]
}
