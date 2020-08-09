// D71 Directory
//
// The directory is always located on track 18 (like the D64), however, on double
// sided disks the extra BAM information is stored on track 53 at sector 0.
// The remaining sectors on track 53 are marked in the BAM as allocated.
//
// Sector 0 of track 18 contains the BAM (Block Availability Map) and disk name/ID,
// while sectors 1-18 contain the 32-byte directory entries.
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
package d71

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mrcook/retroio/commodore/d64"
	"github.com/mrcook/retroio/commodore/disk"
)

const bamExtraTrackNumber = 53 - 1 // tracks count from 1

// Directory represents a D71 directory as stored on a disk.
type Directory struct {
	disk *disk.Disk

	bam         d64.BlockAvailabilityMap
	bamExtra    [35]disk.BamEntry24Bit
	directories []d64.DirectoryFile
}

func newDirectory(dsk *disk.Disk) *Directory {
	return &Directory{disk: dsk}
}

func (d *Directory) Read() error {
	d.bam = d64.BlockAvailabilityMap{}
	if err := d.bam.Read(d.disk.Tracks[d64.DirectoryTrackNumber].Sectors[0]); err != nil {
		return fmt.Errorf("error reading BAM: %w", err)
	}

	if err := d.readExtraBam(); err != nil {
		return fmt.Errorf("error reading Extra BAM: %w", err)
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

	for t := d64.DirectoryTrackNumber; t < len(d.disk.Tracks); {
		sectors := d.disk.Tracks[t].Sectors

		// Read all directories from each allocated sector
		for s := startSector; s < uint8(len(sectors)); {
			reader := bytes.NewReader(sectors[s][:])
			entries := make([]disk.DirectoryFile, d64.DirEntriesPerSector)
			if err := binary.Read(reader, binary.LittleEndian, entries); err != nil {
				return fmt.Errorf("error reading directory entry from sector: %w", err)
			}

			for _, dir := range entries {
				if dir.FileType == 0 && dir.FileSizeInSectors == 0 {
					continue
				}

				entry := d64.DirectoryFile{
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

// In order to manage the extra 35 tracks, most of the extra BAM information is
// stored on track 53/0, and the remaining sectors on track 53 are marked in
// the BAM as allocated.
func (d *Directory) readExtraBam() error {
	sector := d.disk.Tracks[bamExtraTrackNumber].Sectors[0]

	reader := bytes.NewReader(sector[:])

	// The FreeSectors for tracks 36-70 are stored in the main BAM, the
	// remaining 3-bytes of each entry are read from 53/0.
	for i := 0; i < 35; i++ {
		var bitmap [3]uint8
		if _, err := reader.Read(bitmap[:]); err != nil {
			return err
		}

		d.bamExtra[i] = disk.BamEntry24Bit{
			FreeSectors:  d.bam.FreeSectorCount[i],
			SectorBitmap: bitmap,
		}
	}

	return nil
}
