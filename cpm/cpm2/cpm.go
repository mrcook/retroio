// CP/M 2.2 Disc Format
//
// Designed to work with a much larger range of discs than CP/M 1.4.
// The reserved tracks will contain an image of CP/M 2.2, used when the system
// is rebooted. Discs can be formatted as data only discs, in which case they
// have no system tracks and cannot be used to reboot the system.
//
// Reference: http://www.seasip.info/Cpm/format22.html
package cpm2

import "retroio/cpm/cpm1"

// Disk Parameter Header
//
// The DPH retains the CPM 1.4 specification.
type DiskParameterHeader struct {
	cpm1.DiskParameterHeader
}

// Disk Parameter Block
//
// The DPB retains the CPM 1.4 specification.
type DiskParameterBlock struct {
	cpm1.DiskParameterBlock
}

// Disk Directory
//
// The CP/M 2.2 directory is the same header as with the 1.4 standard,
// although there are minor changes to the parameter usage.
//
// The directory is a sequence of directory entries (also called extents),
// which contain 32 bytes of the following structure:
//
// UU F1 F2 F3 F4 F5 F6 F7 F8 T1 T2 T3 EX S1 S2 RC   .FILENAMETYP....
// AL AL AL AL AL AL AL AL AL AL AL AL AL AL AL AL   ................
type Directory struct {
	// UU - User Number
	// Allows multiple files of the same name to coexist on the disc.
	// - 0-15:  Used for file, status is the user number
	// - 16-31: Used only on CP/M 3+, for password extent
	// - 32:    Disc label
	// - 0xE5:  File deleted
	UserNumber uint8

	// Fn - filename
	// Tn - file type extension name
	//
	// May consist of any printable 7 bit ASCII character but: < > . , ; : = ? * [ ].
	// The file name must not be empty, the extension may be empty.
	// Both are padded with blanks.
	//
	// The highest bit of each character of the file extension are used as
	// attributes, and have the following meanings:
	// - T1: top bit is set if the file is read-only.
	// - T2: is set if the file is a system file ("hidden" on other systems).
	// - T3: archived
	//
	// Public files (visible under each user number) are not supported by CP/M 2.2,
	// but there is a patch and some free CP/M clones support them without any patches.
	Filename [8]uint8
	FileType [3]uint8

	// EX - Extent counter, low byte
	ExtentLow uint8 // Has values from 0-31

	// S1 - reserved, set to 0.
	S1 uint8

	// S2 - Extent counter, high byte.
	//
	// An extent is the portion of a file controlled by one directory entry.
	// If a file takes up more blocks than can be listed in one directory entry,
	// it is given multiple entries, distinguished by their EX and S2 bytes.
	// The formula is: Entry number = ((32*S2)+EX) / (exm+1) where exm is the
	// extent mask value from the Disc Parameter Block.
	ExtentHigh uint8

	// RC - Number of records (1 record=128 bytes) used in this extent, low byte.
	// The total number of records used in this extent is:
	//
	// (EX & exm) * 128 + RC
	//
	// If RC is 80h, this extent is full and there may be another one on the disc.
	// File lengths are only saved to the nearest 128 bytes.
	RecordCount uint8

	// AL - Allocation.
	// Each AL is the number of a block on the disc. If an AL number is zero,
	// that section of the file has no storage allocated to it (ie it does not
	// exist). For example, a 3k file might have allocation 5,6,8,0,0.... - the
	// first 1k is in block 5, the second in block 6, the third in block 8.
	// AL numbers can either be 8-bit (if there are fewer than 256 blocks on the
	// disc) or 16-bit (stored low byte first).
	Allocation [16]uint8
}
