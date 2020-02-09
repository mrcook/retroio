package disk

import "fmt"

// Emulated disk storage type
type storageID uint8

const (
	D64 storageID = iota // all variations: 35, 40, and 42 tracks
	D71
	D81
)

// layout is a custom type for describing the attributes of an emulated disk
type layout struct {
	format           storageID
	trackCount       uint8
	totalSectorCount uint16
	fileSize         uint32
	errorBytes       uint16
	description      string
}

// Commodore disk layouts for all emulated formats
var diskLayouts = []layout{
	{D64, 35, 683, 174848, 0, "Standard D64 (1540/41)"},
	{D64, 35, 683, 175531, 683, "Standard D64 (1540/41) with Error Bytes"},
	{D64, 40, 768, 196608, 0, "Standard D64 (1541)"},
	{D64, 40, 768, 197376, 768, "Standard D64 (1541) with Error Bytes"},
	{D64, 42, 802, 205312, 0, "Extended 42 Track D64"},
	{D64, 42, 802, 206114, 802, "Extended 42 Track D64 with Error Bytes"},
	{D71, 70, 1366, 349696, 0, "Standard D71 (1571)"},
	{D71, 70, 1366, 351062, 1366, "Standard D71 (1571) with Error Bytes"},
	{D81, 80, 3200, 819200, 0, "Standard D81 (1581)"},
	{D81, 80, 3200, 822400, 3200, "Standard D81 (1581) with Error Bytes"},
}

func layoutByFileSize(fileSize uint32) (layout, error) {
	for _, v := range diskLayouts {
		if v.fileSize == fileSize {
			return v, nil
		}
	}
	return layout{}, nil
}

// Disk DOS version codes.
// DOS version $41 is the most common, being used for the 1540/41/71 drives,
// and other variants such as SpeedDOS and Dolphin-DOS.
// Any other value enables "soft write protection" (see `docs.md`)
var diskDosVersions = map[uint8]string{
	0x00: "Not Set",
	0x41: "('A) CBM DOS v2.6 (1540/41/71), and others",
	0x44: "('D) CBM DOS v3.x (1581)",
	0x50: "('P') PrologicDOS 1541 and ProSpeed 1571 2.0",
}

// DOS type codes.
// NOTE: the 1581 is possibly CBM DOS v10.0.
var dosTypes = map[string]string{
	"2A": "CBM DOS v2.6 (1540/41/71), and others",
	"2P": "PrologicDOS 1541 and ProSpeed 1571 2.0",
	"3D": "CBM DOS v3.x (1581)",
	"4A": "Professional DOS Release",
}

// Disk geometry for all tracks, for all types of emulated disks
type trackGeometry struct {
	side            uint8
	startTrack      uint8
	endTrack        uint8
	sectorsPerTrack int
	totalSectors    int
}

// Track layouts for all 6 variations of the single-sided D64, and the
// double-sided D71 and D81 formats.
var trackGeometries = map[storageID][]trackGeometry{
	D64: {
		{0, 1, 17, 21, 357},
		{0, 18, 24, 19, 133},
		{0, 25, 30, 18, 108},
		{0, 31, 35, 17, 85},
		{0, 36, 40, 17, 85},
		{0, 41, 42, 17, 34},
	},
	D71: {
		{0, 1, 17, 21, 357},
		{0, 18, 24, 19, 133},
		{0, 25, 30, 18, 108},
		{0, 31, 35, 17, 85},
		{1, 36, 52, 21, 357},
		{1, 53, 59, 19, 133},
		{1, 60, 65, 18, 108},
		{1, 66, 70, 17, 85},
	},
	D81: {
		{0, 1, 40, 40, 1600},
		{1, 41, 80, 40, 1600},
	},
}

func trackLayout(diskType storageID, trackNumber uint8) (trackGeometry, error) {
	for _, t := range trackGeometries[diskType] {
		if trackNumber >= t.startTrack && trackNumber <= t.endTrack {
			return t, nil
		}
	}
	return trackGeometry{}, fmt.Errorf("invalid track number: #%d", trackNumber)
}
