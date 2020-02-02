// Package dsk implements reading Amstrad DSK image files.
//
// Additional DSK geometry documentation can be found in the `docs.md` file.
// Note: all WORD and DWORD values are stored in low/high byte order.
package dsk

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"retroio/amstrad/dsk/amsdos/cat"
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

	AmsDos AmsDos
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

	// Read the contents of the disk as AMSDOS format
	d.AmsDos = AmsDos{}
	if err := d.AmsDos.Read(d); err != nil {
		return errors.Wrap(err, "AMSDOS read error")
	}

	return nil
}

// DisplayGeometry prints the disk, track and sector metadata to the terminal.
func (d DSK) DisplayGeometry() {
	fmt.Println("DISK INFORMATION:")
	fmt.Println(d.Info)

	for _, track := range d.Tracks {
		sectorSize, _ := sectorSizeMap[track.SectorSize]

		str := fmt.Sprintf("SIDE %d, TRACK %02d: ", track.Side, track.Track)
		if track.SectorsCount == 0 {
			str += "[Track is blank]"
		}
		str += fmt.Sprintf("%02d sectors", track.SectorsCount)
		str += fmt.Sprintf(" (%d bytes)", sectorSize)
		if int(track.SectorsCount) != len(track.Sectors) {
			str += fmt.Sprintf(" WARNING only %d sectors read", len(track.Sectors))
		}
		fmt.Println(str)
	}
}

// CommandDir displays the disk directory to the terminal.
func (d DSK) CommandDir() {
	commandCat, err := cat.CommandCat(d.AmsDos.DPB.BlockCount, d.AmsDos.Directories)
	if err != nil {
		fmt.Printf("CAT command error: %s", err)
		return
	}

	fmt.Printf("Drive %c: user %d\n", commandCat.Drive, commandCat.User)
	fmt.Println()

	// Print listing in two columns
	maxRowsLeft, maxRowsRight := recordRowCounts(len(commandCat.Records))
	for i := 0; i < maxRowsLeft; i++ {
		row := commandCat.Records[i].String()
		if i < maxRowsRight {
			row += fmt.Sprintf("   %s", commandCat.Records[maxRowsLeft+i].String())
		}
		fmt.Println(row)
	}

	fmt.Println()
	fmt.Printf("%3dK free\n", commandCat.FreeSpace)

	if commandCat.HiddenFiles > 0 {
		fmt.Println()
		pluralized := ""
		if commandCat.HiddenFiles > 1 {
			pluralized = "s"

		}
		fmt.Printf("* %d hidden file%s\n", commandCat.HiddenFiles, pluralized)
	}
}

func recordRowCounts(records int) (int, int) {
	left := records / 2
	if records%2 > 0 {
		left += 1
	}
	right := records / 2

	return left, right
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
