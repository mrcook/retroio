// Each Track Block comprises a Track Information Block and sector data.
// The sector data is always at an offset of &100 bytes from the start of the track block.
// The data for the next track in the disc image immediately follows the data for the current track.
//
// The first Track Block is located at offset &100 in the disk image file.
// The track block starts with the Track Information Block.
//
// If the track contains different sized sectors, then the data allocated must be the size of the biggest sector.
// The "sector size" parameter is used to calculate the location of the sector data.
//
// * Sector data always follows Track Information Block at offset &100 from the start of the track information block.
// * Sector data is stored in the same order as the sectors in the sector info block.
package dsk

import (
	"github.com/pkg/errors"

	"github.com/mrcook/retroio/librio/storage"
)

// Track information block
//
// * "sector size" parameter is used to calculate the location of each sector's
//    data. Therefore, the data allocated for each sector must be the same.
// * "number of sectors" is used to identify the number of valid entries in the
//   sector information list.
type TrackInformation struct {
	Identifier [13]byte // Identifier: "Track-Info\r\n"
	Unused1    [3]uint8 // unused
	Track      uint8    // track number
	Side       uint8    // side number

	// EXTENDED format extension by John Elliott.
	//
	// Defines the rate at which data was written to the track.
	// This value applies to the entire track.
	//   0: Unknown
	//   1: Single or double density
	//   2: High Density
	//   3: Extended density
	DataRate uint8 // unused for STANDARD disk images

	// EXTENDED format extension by John Elliott.
	//
	// Defining the encoding used to write the data to the disc and the structure
	// of the data on the disc, including the layout of the sectors.
	// This value applies to the entire track.
	//   0: Unknown (as existing files would have zeroes in these bytes)
	//   1: FM
	//   2: MFM
	RecordingMode uint8 // unused for STANDARD disk images

	SectorSize   uint8 // sector size (enum 0-3)
	SectorsCount uint8 // number of sectors
	GapLength    uint8 // GAP#3 length
	FillerByte   uint8 // filler byte

	Sectors    []SectorInformation // Sector Information List
	SectorData []SectorData        // Sector data, starting at 0x0100 from start of track
}

// Read the track information header.
func (t *TrackInformation) Read(reader *storage.Reader) error {
	copy(t.Identifier[:], reader.ReadBytes(13))
	copy(t.Unused1[:], reader.ReadBytes(3))
	t.Track = reader.ReadByte()
	t.Side = reader.ReadByte()
	t.DataRate = reader.ReadByte()
	t.RecordingMode = reader.ReadByte()
	t.SectorSize = reader.ReadByte()
	t.SectorsCount = reader.ReadByte()
	t.GapLength = reader.ReadByte()
	t.FillerByte = reader.ReadByte()

	if err := t.readSectorInformationBlocks(reader); err != nil {
		return err
	}

	if err := t.readSectorData(reader); err != nil {
		return err
	}

	return nil
}

func (t *TrackInformation) readSectorInformationBlocks(reader *storage.Reader) error {
	for i := 0; i < int(t.SectorsCount); i++ {
		sector := SectorInformation{}
		if err := sector.Read(reader); err != nil {
			return errors.Wrapf(err, "error reading sector #%d", i+1)
		}
		t.Sectors = append(t.Sectors, sector)
	}
	return nil
}

func (t *TrackInformation) readSectorData(reader *storage.Reader) error {
	if err := t.setBufferToDataAddress(reader); err != nil {
		return err
	}

	for i, s := range t.Sectors {
		data, err := s.dataRead(reader)
		if err != nil {
			return errors.Wrapf(err, "error reading sector #%d", i)
		}
		t.SectorData = append(t.SectorData, SectorData{Data: data})
	}
	return nil
}

func (t TrackInformation) setBufferToDataAddress(reader *storage.Reader) error {
	blockSize := int(t.SectorsCount) * sectorInformationBlockSize
	usedBytes := trackInformationHeaderSize + blockSize

	_, err := reader.Discard(SectorDataStartAddress - usedBytes)
	if err != nil {
		return errors.Wrap(err, "error moving reader position to 0x0100")
	}

	return nil
}
