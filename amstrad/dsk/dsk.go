// Package dsk implements reading Amstrad DSK image files.
//
// Additional DSK geometry documentation can be found in the `docs.md` file.
// Note: all WORD and DWORD values are stored in low/high byte order.
package dsk

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"retroio/storage"
)

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

	Info   DiskInformation
	Tracks []TrackInformation
}

func New(reader *storage.Reader) *DSK {
	return &DSK{reader: reader}
}

func (d *DSK) Read() error {
	d.Info = DiskInformation{}
	if err := d.Info.Read(d.reader); err != nil {
		return errors.Wrap(err, "error reading the disk information block")
	}

	for i := 0; i < int(d.Info.Tracks); i++ {
		track := TrackInformation{}
		if err := track.Read(d.reader); err != nil {
			return errors.Wrapf(err, "error reading track #%d", i+1)
		}
		d.Tracks = append(d.Tracks, track)
	}

	return nil
}

// DisplayGeometry prints the disk, track and sector metadata to the terminal.
func (d DSK) DisplayGeometry() {
	fmt.Println("DISK INFORMATION:")
	fmt.Println(d.Info)

	for _, track := range d.Tracks {
		str := fmt.Sprintf("SIDE %d, TRACK %02d: ", track.Side, track.Track)
		if track.SectorsCount == 0 {
			str += "[Track is blank]"
		}
		str += fmt.Sprintf("%02d sectors", track.SectorsCount)
		if len(track.Sectors) > 0 {
			str += fmt.Sprintf(" (%d bytes)", track.Sectors[0].SectorByteSize())
		}
		if int(track.SectorsCount) != len(track.Sectors) {
			str += fmt.Sprintf(" WARNING only %d sectors read", len(track.Sectors))
		}
		fmt.Println(str)
	}
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
