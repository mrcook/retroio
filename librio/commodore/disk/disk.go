// Package disk implements reading of Commodore D64, D71, and D81 disk images.
//
// The D64, D71, D81 disk formats were documented by Peter Schepers. These
// specifications can be found at:
// http://vice-emu.sourceforge.net/vice_16.html
// https://ist.uwaterloo.ca/~schepers/formats.html
package disk

import (
	"github.com/mrcook/retroio/librio/storage"
)

// Images for Commodore disks are basically a byte-for-byte copy of the different
// disk variations: 1540, 1541, 1571, and 1581, each with their own sub-variations.
//
// The standard disk formats are:
//   - D64: 35 tracks,  683 sectors, 256 bytes per sector, 174848 byte file size.
//   - D71: 70 tracks, 1366 sectors, 256 bytes per sector, 349696 byte file size.
//   - D81: 80 tracks, 3200 sectors, 256 bytes per sector, 819200 byte file size.
//
// The D64 has two other variations: 40 track (768 sectors) and the rare 42 track (802 sectors).
// Sectors per track varies depending of the location on disk.
// Some variations append errors bytes to the end of the disk; 1 byte per sector.
type Disk struct {
	Tracks []Track // 35 - 80 depending on media type

	// Optional error bytes for some disk variations (see d64.md)
	// - D64: 683, 768, or 802 bytes
	// - D71: 1366 bytes
	// - D81: 3200 bytes
	Errors []uint8

	// Internal usage only
	variation layout
}

// Read the disk image, returning errors if the disk layout (track/sector counts)
// are of unexpected values for the expected media type.
func (d *Disk) Read(reader *storage.Reader) error {
	// var i uint8
	//
	// // NOTE: track counting starts at 1, not 0!
	// for i = 1; i <= d.variation.tracks; i++ {
	// 	geometry, err := trackLayout(d.variation.mediaType, i)
	// 	if err != nil {
	// 		return fmt.Errorf("unable to determine geomery for this track number: %w", err)
	// 	}
	// 	track := Track{Number: i, Geometry: geometry}
	//
	// 	if err := track.Read(reader); err != nil {
	// 		return fmt.Errorf("error reading track #%d: %w", i, err)
	// 	}
	//
	// 	d.Tracks = append(d.Tracks, track)
	// }
	//
	// totalSectorCount := 0
	// for _, track := range d.Tracks {
	// 	totalSectorCount += len(track.Sectors)
	// }
	// if totalSectorCount != int(d.variation.totalSectors) {
	// 	return fmt.Errorf("incorrect sector count for disk, expected %d, got %d\n", d.variation.totalSectors, totalSectorCount)
	// }

	return nil
}
