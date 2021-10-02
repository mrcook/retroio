package trd

// DiskType represents the possible values of the disk type
//
// A TR-DOS disk contains 1 or 2 sides with 40 or 80 tracks per side and with
// 16 sectors per track. Each sector is 256 bytes long.
// There are 4 different disk formats.
type DiskType uint8

const (
	Tracks80Sides2 DiskType = 0x16
	Tracks40Sides2 DiskType = 0x17
	Tracks80Sides1 DiskType = 0x18
	Tracks40Sides1 DiskType = 0x19
)
