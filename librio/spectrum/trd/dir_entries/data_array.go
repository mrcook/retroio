package dir_entries

import "github.com/mrcook/retroio/librio/spectrum/trd/dir_entries/ext"

type DataArray struct {
	// Filename, padded with spaces.
	// NOTE: if the first character is
	//   0x00: then it's the end of the directory.
	//   0x01: indicates a deleted file, which is still present on the disk.
	Filename [8]uint8

	FileExtension   ext.FileExtension // a character that describing the file type
	Unused          uint16            //
	DataArrayLength uint16            // Length of data array
	LengthInSectors uint8             // File length in sectors
	StartingSector  uint8             // Starting sector, calculated as: start track*16+start sector
	StartingTrack   uint8             // Starting (logical) track
}
