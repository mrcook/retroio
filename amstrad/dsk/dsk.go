// Package dsk implements reading of Amstrad DSK image files as specified at:
// http://cpctech.cpc-live.com/docs/dsk.html
//
// * Track 0 (or Track 0 side 0 for double sided disks), if track data exists,
//   will immediately follow the Disc Information Block and will start at
//   offset &100 in the disc image file.
// * All tracks must have a "Track Information Block"
// * Track lengths are stored in the same order as the tracks in the image e.g.
//   In the case of a double sided disk: Track 0 side 0, Track 0 side 1, Track 1 side 0, etc.
// * The track blocks are stored in increasing order 0..number of tracks, with
//   alternating sides interleaved if the disc image describes a double sided
//   disk. e.g. if the disk image represents a double sided disk, the order of tracks is:
//   - track 0 side 0,
//   - track 0 side 1,
//   - track 1 side 0,
//   - track 1 side 1....
//   - track (number of tracks-1) side 0, track (number of tracks-1) side 1
//
// The tracks are always ordered in this way regardless of the disc-format described by the disc image.
//
// A standard disk image can be used to describe a copy-protected disk, but will often result
// in a file which is larger than the same disk described by a extended disk image.
// For a standard disk image to represent a copy-protected disk:
//   - All track sizes in the standard disk image must be the same. This value therefore would
//     be the size of the largest track, and other tracks would have unused space in them.
//   - All sector sizes within each track must be the same size, but not necessarily the same
//     size as the sectors for another track. If a track contained different sized sectors,
//     the size of the largest sector should be used. This would result in some wasted space.
//
// General format:
//
// Single sided DSK images:
// * Disc Information Block
// * Track 0 data
//   - Track Information Block
//   - Sector data
// * Track 1 data
//   - Track Information Block
//   - Sector data
// * . . . .
// * Track (number_of_tracks-1) data
//   - Track Information Block
//   - Sector data
//
// Double sided DSK images:
// * Disc Information Block
// * Track 0 side 0 data
//   - Track Information Block
//   - Sector data
// * Track 0 side 1 data
//   - Track Information Block
//   - Sector data
// * . . . .
// * Track (number_of_tracks-1) side 1 data
//   - Track Information Block
//   - Sector data
//
// Note: all WORD and DWORD values are stored in low/high byte order.
package dsk

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"retroio/storage"
)

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
