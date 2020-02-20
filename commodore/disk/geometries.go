// Each type of Commodore disk media has a different geometry layout; numbers
// of tracks, sectors, etc. That data is initialized here, along with some
// helper functions for accessing that data while reading a disk image.
package disk

import (
	"fmt"

	"retroio/commodore"
)

// layout is a custom type for describing the attributes of an emulated disk
type layout struct {
	mediaType    commodore.MediaType
	diskSize     uint32
	tracks       uint8
	totalSectors uint16
	errorBytes   uint16
	description  string
}

// Disk layouts for all supported commodore disk media
var diskLayouts = []layout{
	{commodore.D64, 174848, 35, 683, 0, "Standard D64 (1540/41)"},
	{commodore.D64, 175531, 35, 683, 683, "Standard D64 (1540/41) with error bytes"},
	{commodore.D64, 196608, 40, 768, 0, "Standard D64 (1541)"},
	{commodore.D64, 197376, 40, 768, 768, "Standard D64 (1541) with error bytes"},
	{commodore.D64, 205312, 42, 802, 0, "Extended 42 Track D64"},
	{commodore.D64, 206114, 42, 802, 802, "Extended 42 Track D64 with error bytes"},
	{commodore.D71, 349696, 70, 1366, 0, "Standard D71 (1571)"},
	{commodore.D71, 351062, 70, 1366, 1366, "Standard D71 (1571) with error bytes"},
	{commodore.D81, 819200, 80, 3200, 0, "Standard D81 (1581)"},
	{commodore.D81, 822400, 80, 3200, 3200, "Standard D81 (1581) with error bytes"},
}

func layoutForMedia(mediaType commodore.MediaType, fileSize uint32) (layout, error) {
	for _, v := range diskLayouts {
		if v.diskSize != fileSize {
			continue
		}

		// TODO: return a user friendly media type
		if v.mediaType != mediaType {
			return layout{}, fmt.Errorf("invalid media type #%d for file size %d", mediaType, fileSize)
		}

		return v, nil
	}

	return layout{}, fmt.Errorf("no disk layout found for media type #%d and file size %d", mediaType, fileSize)
}

// Disk DOS version codes.
// DOS version $41 is the most common, being used for the 1540/41/71 drives,
// and other variants such as SpeedDOS and Dolphin-DOS.
// Any other value enables "soft write protection" (see `docs.md`)
var DosVersions = map[uint8]string{
	0x00: "Not Set",
	0x41: "('A') CBM DOS v2.6 (1540/41/71), and others",
	0x44: "('D') CBM DOS v3.x (1581)",
	0x50: "('P') PrologicDOS 1541 and ProSpeed 1571 2.0",
}

// DOS type codes.
// NOTE: the 1581 is possibly CBM DOS v10.0.
var DosTypes = map[string]string{
	"2A": "CBM DOS v2.6",
	"2P": "PrologicDOS 1541 and ProSpeed 1571 2.0",
	"3D": "CBM DOS v3.x",
	"4A": "Professional DOS Release",
}

// Disk track geometry for all types of emulated disks, as a range of track numbers.
type trackGeometry struct {
	side            uint8 // NOTE: `0` and `1` not `1` and `2`
	startTrack      uint8
	endTrack        uint8
	sectorsPerTrack int
	totalSectors    int // total sectors for this range of tracks
}

// Track layouts for all 6 variations of the single-sided D64, and the
// double-sided D71 and D81 formats. This data is used to aid in reading the
// correct number of track sectors for a given media type.
var trackGeometries = map[commodore.MediaType][]trackGeometry{
	commodore.D64: {
		{0, 1, 17, 21, 357},
		{0, 18, 24, 19, 133},
		{0, 25, 30, 18, 108},
		{0, 31, 35, 17, 85},
		{0, 36, 40, 17, 85},
		{0, 41, 42, 17, 34},
	},
	commodore.D71: {
		{0, 1, 17, 21, 357},
		{0, 18, 24, 19, 133},
		{0, 25, 30, 18, 108},
		{0, 31, 35, 17, 85},
		{1, 36, 52, 21, 357},
		{1, 53, 59, 19, 133},
		{1, 60, 65, 18, 108},
		{1, 66, 70, 17, 85},
	},
	commodore.D81: {
		{0, 1, 40, 40, 1600},
		{1, 41, 80, 40, 1600},
	},
}

// Returns the layout for the track number of a media type.
func trackLayout(diskType commodore.MediaType, trackNumber uint8) (trackGeometry, error) {
	for _, t := range trackGeometries[diskType] {
		if trackNumber >= t.startTrack && trackNumber <= t.endTrack {
			return t, nil
		}
	}
	return trackGeometry{}, fmt.Errorf("invalid track number: #%d", trackNumber)
}
