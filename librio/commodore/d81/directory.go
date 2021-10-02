// D81 Directory
//
// The directory track should be contained totally on track 40. Sectors 1 and 2
// contain the BAM (Block Availability Map), Sectors 3-39 contain the entries.
//
// Sector 0 holds the disk name and ID.
// The first directory sector is always 40/3, even though the t/s pointer
// (first two bytes) at 40/0 might point somewhere else. It goes linearly up
// the sector count, 3-4-5-6-etc. (an interleave of 1).
// Each sector holds up to eight entries.
//
// NOTE: unlike the D64/D71, the interleave for a D81 is `1` for both files and directories.
package d81

// Directory represents a D81 directory as stored on a disk.
// There are a maximum of 296 entries.
type Directory struct {
	header   Header
	bamSide1 BlockAvailabilityMap
	bamSide2 BlockAvailabilityMap

	// TODO: is it okay to initialize with 144 entries?
	// TODO: do we care to handle those rare disks with more?
	directories [144]DirectoryFile
}

func (d *Directory) Read() error {
	// d.header = Header{}
	// if err := d.header.Read(d.disk.Tracks[DirectoryTrackNumber].Sectors[0]); err != nil {
	// 	return fmt.Errorf("error reading Header: %w", err)
	// }
	//
	// d.bamSide1 = BlockAvailabilityMap{}
	// if err := d.bamSide1.Read(d.disk.Tracks[DirectoryTrackNumber].Sectors[1]); err != nil {
	// 	return fmt.Errorf("error reading Side 1 BAM: %w", err)
	// }
	//
	// d.bamSide2 = BlockAvailabilityMap{}
	// if err := d.bamSide2.Read(d.disk.Tracks[DirectoryTrackNumber].Sectors[2]); err != nil {
	// 	return fmt.Errorf("error reading Side 2 BAM: %w", err)
	// }
	//
	// if err := d.readDirectory(); err != nil {
	// 	return fmt.Errorf("directory track error: %w", err)
	// }

	return nil
}

// Read all 32-byte directory entries, following the T/S chain.
func (d *Directory) readDirectory() error {
	// Header is in sector[0], the two BAMs are in sectors 1 and 2,
	// the dirs start at sector #3
	// var startSector uint8 = 3

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
