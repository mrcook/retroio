// Package dsk implements reading Amstrad DSK image files.
//
// Additional DSK geometry documentation can be found in the `d64.md` file.
// Note: all WORD and DWORD values are stored in low/high byte order.
package dsk

import (
	"strings"

	"github.com/mrcook/retroio/librio/amstrad/amsdos"
	"github.com/mrcook/retroio/librio/storage"
	"github.com/pkg/errors"
)

const (
	trackInformationHeaderSize = 24
	sectorInformationBlockSize = 0x08   // bytes
	SectorDataStartAddress     = 0x0100 // 256-bytes
)

// sectorSizeMap maps the track/sector size enum byte to its size in bytes value.
var SectorSizeMap = map[uint8]uint16{
	0: 128,
	1: 256,
	2: 512,
	3: 1024,
}

// DSK image format
//
// Track data (if it exists) will immediately follow the Disc Information
// Block, with track #0 starting at offset 0x0100 in the image file.
// Single sided disk tracks are stored sequentially.
// Double sided disk track order is:
//   track 0 side 0
//   track 0 side 1
//   track 1 side 0
//   track 1 side 1
//   etc.
// NOTE: tracks are always ordered in this way regardless of the disc format
// described by the disc image.
type DSK struct {
	reader *storage.Reader

	Info   DiskInformation // 256 bytes long
	Tracks []TrackInformation

	AmsDos amsdos.AmsDos
}

func (d *DSK) Read() error {
	d.Info = DiskInformation{}

	if err := d.Info.Read(d.reader); err != nil {
		return errors.Wrap(err, "error reading the disk information block")
	}

	for i := 0; i < int(d.Info.Tracks); i++ {
		track := TrackInformation{}

		if d.Info.isStandardDisk() || (d.Info.isExtendedDisk() && d.Info.TrackSizeTable[i] != 0) {
			if err := track.Read(d.reader); err != nil {
				return errors.Wrapf(err, "error reading track #%d", i+1)
			}
		}

		d.Tracks = append(d.Tracks, track)
	}

	// Read the contents of the disk as AMSDOS format
	d.AmsDos = amsdos.AmsDos{}
	if err := d.AmsDos.Read(d); err != nil {
		return errors.Wrap(err, "AMSDOS read error")
	}

	return nil
}

func reformatIdentifier(identifier []byte) string {
	var idBytes []byte
	for _, b := range identifier {
		if b > 0 {
			idBytes = append(idBytes, b)
		}
	}

	id := strings.Trim(string(idBytes), "\r\n")
	parts := strings.Split(id, "\r\n")

	return strings.Join(parts, ", ")
}
