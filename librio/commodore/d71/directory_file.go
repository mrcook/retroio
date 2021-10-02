package d71

import "github.com/mrcook/retroio/librio/commodore/disk"

// DirectoryFile contains metadata for each entry, along with the 32-byte
// directory file entry found on the disk.
type DirectoryFile struct {
	Track  uint8
	Sector uint8

	Filename string
	FileType disk.FileType

	DirEntry disk.DirectoryFile
}
