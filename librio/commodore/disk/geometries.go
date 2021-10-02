// Each type of Commodore disk media has a different geometry layout; numbers
// of tracks, sectors, etc. That data is initialized here, along with some
// helper functions for accessing that data while reading a disk image.
package disk

import (
	"fmt"

	"github.com/mrcook/retroio/commodore"
)

// Custom type for the different storage mediums
type MediaType uint8

const (
	Unknown MediaType = iota
	T64
	TAP
	D64 // all variations: 35, 40, and 42 tracks
	D71
	D81
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
