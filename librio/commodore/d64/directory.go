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
	"github.com/mrcook/retroio/librio/storage"
)

// Directory represents a D64 directory as stored on a disk.
// There are 144 max entries, unless it is a multi-track directory.
type Directory struct {
	bam         BlockAvailabilityMap
	directories []DirectoryFile
}

func (d *Directory) Read(reader *storage.Reader) error {
	// d.bam = BlockAvailabilityMap{}
	// if err := d.bam.Read(d.disk.Tracks[DirectoryTrackNumber].Sectors[0]); err != nil {
	// 	return fmt.Errorf("error reading BAM: %w", err)
	// }
	//
	// if err := d.readDirectory(); err != nil {
	// 	return fmt.Errorf("directory track error: %w", err)
	// }

	return nil
}

// Read all 32-byte directory entries, following the T/S chain.
func (d *Directory) readDirectory() error {
	// BAM is in sector[0], dirs start at sector #1
	// var startSector uint8 = 1

	// for t := DirectoryTrackNumber; t < len(d.disk.Tracks); {
	// 	sectors := d.disk.Tracks[t].Sectors
	//
	// 	// Read all directories from each allocated sector
	// 	for s := startSector; s < uint8(len(sectors)); {
	// 		reader := bytes.NewReader(sectors[s][:])
	// 		entries := make([]disk.DirectoryFile, DirEntriesPerSector)
	// 		if err := binary.Read(reader, binary.LittleEndian, entries); err != nil {
	// 			return fmt.Errorf("error reading directory entry from sector: %w", err)
	// 		}
	//
	// 		for _, dir := range entries {
	// 			if dir.FileType == 0 && dir.FileSizeInSectors == 0 {
	// 				continue
	// 			}
	//
	// 			entry := DirectoryFile{
	// 				Track:    uint8(t),
	// 				Sector:   s,
	// 				Filename: dir.PrintableFilename(),
	// 				FileType: dir.FileTypeFromID(),
	// 				DirEntry: dir,
	// 			}
	// 			d.directories = append(d.directories, entry)
	// 		}
	//
	// 		// Fetch the next sector to jump to
	// 		s = entries[0].NextSector
	//
	// 		// that was the last sector, stop reading the directory
	// 		if s == 0xFF {
	// 			return nil
	// 		}
	//
	// 		nextTrack := int(entries[0].NextTrack)
	// 		if nextTrack == 0x00 {
	// 			return nil
	// 		}
	//
	// 		// If a multi-track directory disk, change tracks
	// 		if t != nextTrack-1 {
	// 			startSector = s
	// 			t = nextTrack - 1
	// 			break // break the sector reading loop
	// 		}
	// 	}
	// }

	return nil
}
