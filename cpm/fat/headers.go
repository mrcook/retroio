// The FAT filesystem in 16-bit CP/M-86
//
// Concurrent CP/M 3.1 and later, and single-user CP/M-86 with BDOS 3.3 and
// later (including DOS Plus), allow CP/M programs to access DOS-formatted
// discs via conventional BDOS calls, emulating (as far as possible) the
// behaviour of a normal CP/M filesystem. The behaviour is probably a good
// starting point for anyone writing a CP/M emulator which uses a hierarchical
// or non-CP/M filesystem.
//
// Test for DOS media
// The test used by the DOS Plus utilities is to perform a Search First
// (BDOS function 11h) with the first 12 bytes of the FCB containing '?' marks.
// On a DOS disc, this will return a "disc label" directory entry, and byte
// 0Fh of this entry will be 80h. On a CP/M disc, it may or may not return the
// label, but if it does byte 0Fh of the label entry will be zero.
//
// Reference: http://www.seasip.info/Cpm/dosplus_fat.html
package fat

import "retroio/cpm/cpm3"

// XDPHS Extended Disk Parameter Headers
//
// An Extended Disk Parameter Header (XDPH) consists of a prefix and a regular
// Disk Parameter Header as described in Section 3. The label of a XDPH
// references the start of the DPH. The fields of the prefix are located at
// relative offsets from the XDPH label.
//
// The XDPHs for each unit of a controller are the only entry points in a
// particular disk drive module. They contain both the DPH for the drive and
// the addresses of the various action routines for that drive, including READ,
// WRITE, and initialization.
type DiskParameterHeader struct {
	// WRITE
	// The WRITE word contains the address of the sector WRITE routine for the drive.
	WriteAddress uint16

	// READ
	// The READ word contains the address of the sector READ routine for the drive.
	ReadAddress uint16

	// LOGIN
	// The LOGIN word contains the address of the LOGIN routine for the drive.
	LoginAddress uint16

	// INIT
	// The INIT word contains the address of the f irst-time initialization
	// code for the drive.
	InitAddress uint16

	// UNIT
	// The UNIT byte contains the drive code relative to the disk controller.
	// This is the value placed in @RDRV prior to calling the READ, WRITE,
	// and LOGIN entry points of the drive.
	DriveUnitCode uint8

	// TYPE
	// The TYPE byte is unused by the BIOS root, and is reserved for the driver
	// to keep the current density or media type to support multiple-format disk
	// subsystems.
	DriveType uint8

	// The remaining fields of the XDPH comprise a standard DPH.
	cpm3.DiskParameterHeader
}

// XIOS Extended DPB
//
// The XIOS returns an EDPB for DOS media.
type DiscParameterBlock struct {
	// extflag
	// On a DOS disc, this will be -1 for FAT12 and -2 for FAT16 (CCP/M-86 3.1 only supports FAT12).
	// On a CP/M-86 disc, a normal DPB is returned, and this will be its SPT (Sectors per track) entry.
	ExtFlag uint16

	// nfats - Number of FATs (normally 2)
	FatCount uint16

	// nfatrecs - The number of physical sectors per FAT.
	PhysicalSectorsPerFat uint16

	// nclstrs - The number of clusters on the disc.
	ClusterCount uint16

	// clsize - The number of bytes per cluster.
	BytesPerCluster uint16

	// fatadd - The number of the sector holding the start of the first FAT.
	FatAdd uint16

	// Now follows the CP/M DPB

	SectorsPerTrack   uint16 // spt ; Number of 128-byte records per track
	BlockShift        uint8  // bsh ; Block shift. 3 => 1k, 4 => 2k, 5 => 4k....
	BlockMask         uint8  // blm ; Block mask. 7 => 1k, 0Fh => 2k, 1Fh => 4k...
	ExtentMask        uint8  // exm ; Extent mask
	BlockCount        uint16 // dsm ; (no. of blocks on the disc)-1
	DirectoryCount    uint16 // drm ; (no. of directory entries)-1
	AllocationBitmap0 uint8  // al0 ; Not used for DOS media
	AllocationBitmap1 uint8  // al1 ; Not used for DOS media
	Checksum          uint16 // cks ; Checksum vector size, 0 or 8000h for a fixed disc, else number root directory entries/4, rounded up.
	Offset            uint16 // off ; Offset, number of reserved tracks

	// 0 => 128-byte sectors
	// 1 => 256-byte sectors
	// 2 => 512-byte sectors
	// ...
	PhysicalShift uint8 // psh ; Physical sector shift
	PhysicalMask  uint8 // phm ; Physical sector mask
}

// DOS File Control Block
// Reference: http://www.seasip.info/Cpm/fcb.html
//
// The File Control Block is a 36-byte data structure (33 bytes in CP/M 1).
//
// NOTE: the DOS control behave as the normal CP/M FCB, but the Allocation area
// (16 bytes at FCB+10h) is replaced with different information.
//
// DR F1 F2 F3 F4 F5 F6 F7 F8 T1 T2 T3 EX S1 S2 RC  .FILENAMETYP...
// FL FL OF 00 CC CC CI CI DD DD CL CL LE LE LE LE  ................
// CR R0 R1 R2                                      ....
type FileControlBlock struct {
	// DR - Drive.
	//     0: default
	//  1-16: for A-P
	//
	// In DOSPLUS, bit 7 can be set to indicate that the operation should work
	// with sub-directories rather than files.
	Drive uint8

	// Fn - Filename, 7-bit ASCII.
	// The top bits of the filename bytes (usually referred to as F1' to F8')
	// have the following meanings:
	// F1'-F4' - User-defined attributes.
	//           Any program can use them in any way it likes. The filename in
	//           the disc directory has the corresponding bits set.
	// F5'-F8' - Interface attributes.
	//           They modify the behaviour of various BDOS functions or
	//           indicate error conditions. In the directory these bits are
	//           always zero.
	Filename [8]uint8

	// Tn - Filetype, 7-bit ASCII.
	// T1' - Read-Only.
	// T2' - System (hidden). System files in user 0 can be opened from other user areas.
	// T3' - Archive. Set if the file has not been changed since it was last copied.
	FileType [3]uint8

	// EX - Set this to 0 when opening a file and then leave it to CP/M.
	// You can rewind a file by setting EX, RC, S2 and CR to 0.
	Extent uint8

	// S1 - Reserved.
	S1 uint8

	// S2 - Reserved.
	S2 uint8

	// RC - Set this to 0 when opening a file and then leave it to CP/M.
	RecordCount uint8

	// FL = Flags:
	//   d7    always set
	//   d6    set if file is a directory
	//   d5    set if file is hidden
	//   d4-d0 are high bits of file number
	Flags [2]uint8

	// OF = File Number, low 8 bits.
	// Number is 0 for 1st dir entry, 1 for 2nd, etc.
	FileNumber uint16

	// CC = current cluster, 0 if none
	CurrentCluster uint16

	// CI = index of current cluster (0, 1, 2, ...) in file chain.
	// Only meaningful if CC is non-zero.
	CurrentClusterIndex uint16

	// DD = DOS date
	DosDate uint16

	// CL = DOS cluster
	DosCluster uint16

	// LE = DOS length
	DosLength uint32

	// CR - Current record within extent.
	// It is usually best to set this to 0 immediately after a file has been
	// opened and then ignore it.
	CurrentRecord uint8

	// Rn - Random access record number (not CP/M 1).
	// A 16-bit value in CP/M 2 (with R2 used for overflow); an 18-bit value in CP/M 3.
	RndRecordNumber [3]uint8
}

// Fake Directory Record
//
// Functions 11h and 12h will return a fake directory record with:
//
// Entry 0 holding the file data:
// UU F1 F2 F3 F4 F5 F6 F7 F8 T1 T2 T3 EX S1 S2 RC .FILENAMETYP....
// FL FL OF 00 00 00 TT TT DD DD CL CL LE LE LE LE ................
//
// Entries 1 & 2 blank
// Entry 3 holding SFCB
type FakeDirectory struct {
	// UU = current CP/M user number, or 20h for volume label.
	UserNumber uint8

	// Fn - Filename, 7-bit ASCII.
	// The top bits of the filename bytes (usually referred to as F1' to F8')
	// have the following meanings:
	// F1'-F4' - User-defined attributes.
	//           Any program can use them in any way it likes. The filename in
	//           the disc directory has the corresponding bits set.
	// F5'-F8' - Interface attributes.
	//           They modify the behaviour of various BDOS functions or
	//           indicate error conditions. In the directory these bits are
	//           always zero.
	Filename [8]uint8

	// Tn - Filetype, 7-bit ASCII.
	// T1' - Read-Only.
	// T2' - System (hidden). System files in user 0 can be opened from other user areas.
	// T3' - Archive. Set if the file has not been changed since it was last copied.
	FileType [3]uint8

	// EX - simulated extent number, 0-31.
	Extent uint8

	// S1 = Last record byte count.
	S1 uint8

	// S2 = simulated extent number, high bits.
	S2 uint8

	// RC - simulated record count.
	RecordCount uint8

	// FL = Flags:
	//   d7    always set, thus "block" number is either >80h (8bit blocks) or >8000h (16bit blocks)
	//   d6    set if file is a directory
	//         only shows up with user "?" in search FCB, or if bit 7 is set in FCB "drive" byte
	//   d5    set if file is hidden
	//   d4-d0 are high bits of file number
	Flags [2]uint8

	// OF = File Number, low 8 bits.
	// Number is 0 for 1st dir entry, 1 for 2nd, etc.
	FileNumber uint8

	Unused [3]uint8 // Really?

	// TT = DOS time
	DosTime uint16

	// DD = DOS date
	DosDate uint16

	// CL = DOS cluster
	DosCluster uint16

	// LE = DOS length
	DosLength uint32
}
