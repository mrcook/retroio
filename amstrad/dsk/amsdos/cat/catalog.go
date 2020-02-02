package cat

import (
	"errors"
	"fmt"
	"sort"

	"retroio/amstrad/dsk/amsdos"
)

// COMMAND: CAT
// Catalogs the disc. Generates a list, in alpha-numeric order, the full names
// of all files found, together with each file's length (to the nearest higher Kbyte).
// The free space left on the disc is also displayed, together with Drive and
// User identification.
func CommandCat(diskMaxBlocks uint16, directories []amsdos.Directory) (*catalog, error) {
	if len(directories) == 0 {
		return nil, errors.New("no directories found")
	}

	cat := &catalog{
		Drive:     'A',
		User:      directories[0].UserNumber,
		FreeSpace: diskMaxBlocks,
		Records:   make([]directoryRecord, 0),
	}

	// TODO: research a proper solution for this free space adjustment
	wasExtent := false

	var lastFilename [8]byte
	var lastFileType [3]byte

	for _, d := range directories {
		if !cat.validDirRecord(&d) {
			continue
		}

		record := newDirectoryRecord(d.Filename, d.FileType, cat.blockCount(d.Allocation))

		cat.FreeSpace -= record.RecordCount

		if lastFilename == d.Filename && lastFileType == d.FileType {
			// add this record count to the the last record
			cat.Records[len(cat.Records)-1].RecordCount += record.RecordCount
			wasExtent = true
		} else {
			if record.Hidden {
				cat.HiddenFiles += 1
			}
			cat.Records = append(cat.Records, record)
		}

		lastFilename = d.Filename
		lastFileType = d.FileType
	}

	// without this adjustment the free space is off by 1.
	if wasExtent {
		cat.FreeSpace -= 1
	}

	cat.alphabetize()

	return cat, nil
}

type catalog struct {
	Drive       byte
	User        uint8
	FreeSpace   uint16
	HiddenFiles int
	Records     []directoryRecord
}

// TODO: is there a better way of checking for valid directory entries?
// Check if a directory entry is actually a directory entry
func (c catalog) validDirRecord(dir *amsdos.Directory) bool {
	// Nemesis (1986) is an example game where the first three checks are required.
	// Roland in the Caves (1984) is an example game where all four checks are required.
	if dir.UserNumber <= 32 && dir.RecordCount <= 0x80 && dir.S1 == 0 && dir.ExtentHigh == 0 {
		return true
	}
	return false
}

// Count the allocation blocks
func (c catalog) blockCount(allocation [16]uint8) uint16 {
	var blocks uint16
	for _, b := range allocation {
		if b > 0 {
			blocks += 1
		}
	}
	return blocks
}

// Sorts the catalog records in ascending order using the filename and file extension
func (c *catalog) alphabetize() {
	sort.Slice(c.Records, func(i, j int) bool {
		if c.Records[i].Filename == c.Records[j].Filename {
			return c.Records[i].FileType < c.Records[j].FileType
		}
		return c.Records[i].Filename < c.Records[j].Filename
	})
}

// directoryRecord is the displayable data for a directory record.
// This is similar to the CP/M Directory, except each entry merges all record extents.
type directoryRecord struct {
	Filename    string
	FileType    string
	RecordCount uint16 // Total record count for all extents of a record

	ReadOnly bool
	Hidden   bool
	Archived bool
}

// Returns a displayable directory record from the given disk entry
func newDirectoryRecord(filename [8]byte, fileType [3]byte, blockCount uint16) directoryRecord {
	record := directoryRecord{
		Filename:    string(filename[:]),
		RecordCount: blockCount,
	}

	// Check file type attributes and clear them
	if record.bitSet(fileType[0], 7) {
		fileType[0] = record.clearBit(fileType[0], 7)
		record.ReadOnly = true
	}
	if record.bitSet(fileType[1], 7) {
		fileType[1] = record.clearBit(fileType[1], 7)
		record.Hidden = true
	}
	if record.bitSet(fileType[2], 7) {
		fileType[2] = record.clearBit(fileType[2], 7)
		record.Archived = true
	}

	record.FileType = string(fileType[:])

	return record
}

// String formatted as an Amstrad CAT listing
// Adds a custom "hidden" marker, although not present on the original Amstrad CAT.
func (d directoryRecord) String() string {
	marker := " "
	if d.Hidden {
		marker = "*"
	}

	return fmt.Sprintf("%s.%s %3dK%s", d.Filename, d.FileType, d.RecordCount, marker)
}

func (d directoryRecord) bitSet(n uint8, pos uint8) bool {
	val := n & (1 << pos)
	return val > 0
}

func (d directoryRecord) clearBit(n uint8, pos uint8) uint8 {
	var mask uint8 = ^(1 << pos)
	return n & mask
}
