// D64 Directory
//
// The directory track of a D64 image is located on track 18.
//
// Sector 0 of track 18 contains the BAM (Block Availability Map) and disk name/ID,
// while sectors 1-18 contain the 32-byte directory file entries.
// There are a maximum of 144 directory entries available:
//   18 sectors * 8 entries per sector = 144 total entries.
//
// The first directory sector is always 18/1, even though the T/S pointer at
// 18/0 (first two bytes) might point somewhere else. A sector interleave of 3
// is used for directories, with the chain links going 18/1, 18/4, 18/7 etc.
//
// Note: it is possible, however unlikely, to extend the directory beyond
// track 18, (some disks do exist like this) but only during reading.
// See `docs.md` for more information on this.
package d64

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"retroio/commodore/disk"
)

const (
	DirectoryTrackNumber = 18 - 1 // tracks count from 1
	DirEntriesPerSector  = 8
)

// Directory represents a D64 directory as stored on a disk.
// There are 144 max entries, unless it is a multi-track directory.
type Directory struct {
	disk *disk.Disk

	bam         BlockAvailabilityMap
	directories []DirectoryFile
}

func newDirectory(dsk *disk.Disk) *Directory {
	return &Directory{disk: dsk}
}

// DirectoryFile contains metadata for each entry, along with the 32-byte
// directory file entry found on the disk.
type DirectoryFile struct {
	Track  uint8
	Sector uint8

	Filename string
	FileType disk.FileType

	DirEntry disk.DirectoryFile
}

func (d *Directory) Read() error {
	d.bam = BlockAvailabilityMap{}
	if err := d.bam.Read(d.disk.Tracks[DirectoryTrackNumber].Sectors[0]); err != nil {
		return fmt.Errorf("error reading BAM: %w", err)
	}

	if err := d.readDirectory(); err != nil {
		return fmt.Errorf("directory track error: %w", err)
	}

	return nil
}

// Read all 32-byte directory entries, following the T/S chain.
func (d *Directory) readDirectory() error {
	// BAM is in sector[0], dirs start at sector #1
	var startSector uint8 = 1

	for t := DirectoryTrackNumber; t < len(d.disk.Tracks); {
		sectors := d.disk.Tracks[t].Sectors

		// Read all directories from each allocated sector
		for s := startSector; s < uint8(len(sectors)); {
			reader := bytes.NewReader(sectors[s][:])
			entries := make([]disk.DirectoryFile, DirEntriesPerSector)
			if err := binary.Read(reader, binary.LittleEndian, entries); err != nil {
				return fmt.Errorf("error reading directory entry from sector: %w", err)
			}

			for _, dir := range entries {
				if dir.FileType == 0 && dir.FileSizeInSectors == 0 {
					continue
				}

				entry := DirectoryFile{
					Track:    uint8(t),
					Sector:   s,
					Filename: dir.PrintableFilename(),
					FileType: dir.FileTypeFromID(),
					DirEntry: dir,
				}
				d.directories = append(d.directories, entry)
			}

			// Fetch the next sector to jump to
			s = entries[0].NextSector

			// that was the last sector, stop reading the directory
			if s == 0xFF {
				return nil
			}

			nextTrack := int(entries[0].NextTrack)
			if nextTrack == 0x00 {
				return nil
			}

			// If a multi-track directory disk, change tracks
			if t != nextTrack-1 {
				startSector = s
				t = nextTrack - 1
				break // break the sector reading loop
			}
		}
	}

	return nil
}

// BAM Layout for the D64 and D71 disks.
//
// D64: BAM information is stored on track 18/0
// D71: BAM is different on the 1571 as it must now manage the 35 extra tracks.
//      Most of the extra BAM information is stored in sector 0 on track 53
//      (the equivalent of track 18, but on the flip side of the disk), with
//      the remaining sectors on track 53 marked as allocated.
//      This makes the remainder of track 53 unusable.
type BlockAvailabilityMap struct {
	// Track/Sector location of the first directory sector - should be set to
	// 18/1 but it doesn't matter, and don't trust what is there, always go
	// to 18/1 for first directory entry.
	FirstDirTrack  uint8
	FirstDirSector uint8

	// Disk DOS version (see `docs.md` for more information)
	DiskDosVersion byte // Usually 'A' ($41)

	// D64: Unused ($00)
	// D71: Double-sided flag
	//      $00 - Single sided disk
	//      $80 - Double sided disk
	NumberOfSidesFlag uint8

	// BAM entries for each track
	// In groups of four bytes per track, starting on track 1.
	Entries [35]disk.BamEntry24Bit

	// Disk Name (padded with $A0)
	DiskName [16]byte

	Filler1 [2]uint8 // Filled with $A0

	DiskID [2]byte

	Unknown uint8 // Usually $A0

	// see `docs.md` for more information
	DosVersion  byte // "2"
	DiskVersion byte // "A"

	Filler2 [4]uint8 // Filled with $A0

	// The following 85 bytes are unused on the standard D64 disks,
	// and are filled with $00 values.
	Unused2 uint8

	// DOLPHIN DOS: track 36-40 BAM entries
	DolphinDosBAM [20]uint8

	// SPEED DOS: track 36-40 BAM entries
	SpeedDosBAM [20]uint8

	Unused3 [9]uint8

	// D71 ONLY
	// Free sector count for the 35 tracks on side 2 of a 1571 disk.
	// These 1-byte per track "free sector" entries for tracks 36-70 are
	// likely included here due to memory restrictions in the 1571 drive.
	// NOTE: in a D64 these values would normally be with the 4-byte BAM entry,
	// but the rest of the entry is contained on 53/0.
	FreeSectorCount [35]uint8
}

func (b *BlockAvailabilityMap) Read(sector disk.Sector) error {
	reader := bytes.NewReader(sector[:])
	if err := binary.Read(reader, binary.LittleEndian, b); err != nil {
		return err
	}

	return nil
}

func (b BlockAvailabilityMap) PrintableDiskName() []byte {
	name := b.DiskName
	for i, c := range name {
		if c == 0xA0 {
			name[i] = ' '
		}
	}
	return name[:]
}
