// CP/M 1.4 disc formats
// Designed to work with 8" 250k discs.
package cpm1

// Disk Parameter Header (Table Format) - used by the BIOS/BDOS
//
// In general, each disk drive has an associated (16-byte) disk parameter
// header which both contains information about the disk drive and provides
// a scratchpad area for certain BIOS operations.
type DiskParameterHeader struct {
	// XLT - Logical Translation Vector Offset
	//
	// Offset of the logical to physical translation vector, if used for this
	// particular drive, or the value 0000H if no sector translation takes
	// place (i.e. the physical and logical sector numbers are the same)
	// Disk drives with identical sector skew factors share the same translation tables.
	TranslationVectorOffset uint16

	// Scratchpad values for use within the BDOS (initial value is unimportant).
	ScratchPad [3]uint16

	// DIRBUF - Directory Buffer
	//
	// Offset of a 128-byte scratchpad area for directory operations within BDOS.
	// All DPH's address the same scratchpad area.
	DirectoryBuffer uint16

	// DPB - Disk Parameter Block Offset
	//
	// Offset of a disk parameter block for this drive. Drives with identical disk
	// characteristics address the same disk parameter block.
	DpbOffset uint16

	// CSV - Change disk scratchpad offset
	//
	// Offset of a scratchpad area used for software check for changed disks.
	// This offset is different for each DPH.
	//
	// The size of the area addressed by CSV is CKS bytes, which is sufficient
	// to hold the directory check information for this particular drive.
	// If CKS = (DRM+1)/4, then you must reserve this number of bytes for the
	// directory check use. If CKS = 0, then no storage is required.
	ChangeDiskOffset uint16

	// ALV - Allocation information offset
	//
	// Offset of a scratchpad are used by the BDOS to keep disk storage
	// allocation information. This offset is different for each DPH.
	//
	// The size of the area addressed by ALV is determined by the maximum
	// number of data blocks allowed for this particular disk, and is
	// computed as (DSM/8)+1.
	AllocationOffset uint16
}

// Disk Parameter Block
//
// The disc statistics are stored in a parameter block (the DPB),
// which contains the following information:
//
// NOTE: The DPB is not stored on disc.
type DiskParameterBlock struct {
	// SPT - Total number of 128-byte logical records per track.
	RecordsPerTrack uint16 // usually 8, 16, or 32

	// BLS is (implicitly) determined by the values of BSH and BLM.
	//
	// BSH - Block shift factor
	// The data allocation block shift factor, determine by the data block
	// allocation size.
	BlockShift uint8
	// BLM - Block mask
	// The block mask which is also determined by the data block allocation size.
	BlockMask uint8

	// EXM - Extent mask
	// The extent mask, determined by the data block allocation size and
	// the number of disk blocks.
	//
	// The value of EXM depends upon both the BLS and whether the DSM value
	// is less than 256 or grater than 255, as shown in the following table:
	//
	//    Maximum EXM values
	//  BLS  | DSM<256 | DSM>255
	//  1024 |    0    |   N/A
	//  2048 |    1    |    0
	//  4096 |    3    |    1
	//  8192 |    7    |    3
	// 16384 |   15    |    7
	ExtentMask uint8

	// DSM - Storage Capacity
	// Determines the total storage capacity of the disk drive.
	// Number of blocks on the disc - 1.
	//
	// The value of DSM is the maximum data block number supported by this
	// particular drive, measured in BLS units. The product BLS times (DSM+1)
	// is the total number of bytes held by the drive and, of course, must be
	// within the capacity of the physical disk, not counting the reserved
	// operating system tracks.
	// DSM must be less than or equal to 7FFFH.
	BlockCount uint16

	// DRM - Directory Count
	// Determines the total number of 32-byte directory entries which can be
	// stored on this drive.
	//
	// The DRM entry is one less than the total number of 32-byte directory
	// entries (DRM-1), and is a 16-bit value. DRM must be less than or equal
	// to (BLS/32 * 16) - 1. DRM determines the values of AL0 and AL1.
	DirectoryCount uint16

	// AL0,AL1 - Allocation Bitmap
	//
	// Determines the reserved directory blocks.
	//
	// The values of AL0 and AL1 are determined by DRM. The two values AL0 and
	// AL1 can together be considered a string of 16-bits, as shown below:
	//     AL0      AL1
	//   01234567 89ABCDEF
	//   11110000 00000000
	//
	// (In this example, the first 4 blocks of the disc contain the directory.)
	//
	// Position 0x00 corresponds to the high order bit of the byte labeled AL0,
	// and 15 (0xF) corresponds to the low order but of the byte labeled AL1.
	// Each bit position reserves a data block for a number of directory
	// entries, thus allowing a total of 16 data blocks to be assigned for
	// directory entries, (bits are assigned starting at 0x00 and filled to the
	// right until position 0x0F). Each directory entry occupies 32 bytes, as
	// shown in the following table:
	//
	//  BLS  | Directory Entries
	//  1024 |   32 * number of bits
	//  2048 |   64 * number of bits
	//  4096 |  128 * number of bits
	//  8192 |  256 * number of bits
	// 16384 |  512 * number of bits
	//
	// Thus, if DRM = 127 (128 directory entries), and BLS = 1024, then there
	// are 32 directory entries per block, requiring 4 reserved blocks. In this
	// case, the 4 high order bits of AL0 are set, resulting in the values
	// AL0 = 0F0H and AL1 = 00H.
	AllocationBitmap0 uint8 // first byte
	AllocationBitmap1 uint8 // second byte

	// CKS ; Checksum
	// Size of the directory check vector.
	//
	// Fixed Media: CKS = 0 (no directory records are checked)
	// Removable:   CSK = (DRM+1)/4, rounded up (DRM = last dir entry number)
	Checksum uint16

	// OFF - Reserved tracks count offset
	//
	// Number of reserved tracks at the beginning of the (logical) disk.
	//
	// Determines the number of tracks which are skipped at the beginning of
	// the physical disk. This value is automatically added whenever `SETTRK`
	// is called, and can be used as a mechanism for skipping reserved
	// operating system tracks, or for partitioning a large disk into smalled
	// segmented sections.
	ReservedTracksOffset uint16
}

func (d *DiskParameterBlock) SetAllocationBitmap(reservedBlocks int) {
	var allocation uint16

	// turn on the bits from left to right, per directory
	for i := 0; i < reservedBlocks; i++ {
		allocation |= 1 << (15 - i)
	}

	d.AllocationBitmap0 = uint8((allocation & 0xFF00) >> 8)
	d.AllocationBitmap1 = uint8(allocation & 0x00FF)
}

// BLS Table
//
// The values of BSH and BLM determine (implicitly) the data allocation
// size BLS, which is not an entry in the DPB.
type blsRcord map[uint16]struct {
	Dirs uint16 // Reserved directory blocks
	BSH  uint8  // Block shift factor
	BLM  uint8  // Block mask
}

var BlsTable = blsRcord{
	1024:  {Dirs: 32, BSH: 3, BLM: 7},
	2048:  {Dirs: 64, BSH: 4, BLM: 15},   // 0x0F
	4096:  {Dirs: 128, BSH: 5, BLM: 31},  // 0x1F
	8192:  {Dirs: 256, BSH: 6, BLM: 63},  // 0x3F
	16384: {Dirs: 512, BSH: 7, BLM: 127}, // 0x7F
}

// Disk Directory
//
// The directory is a sequence of directory entries (also called extents),
// which contain 32 bytes of the following structure:
//
// SS F1 F2 F3 F4 F5 F6 F7 F8 T1 T2 T3 EX S1 S2 RC   .FILENAMETYP....
// AL AL AL AL AL AL AL AL AL AL AL AL AL AL AL AL   ................
type Directory struct {
	// SS = Status.
	//   0  => File exists
	// 0E5h => File deleted
	//  80h => File exists and is hidden. This feature was undocumented
	//         and does not exist in later versions of CP/M.
	Status uint8

	// Fn - filename
	Filename [8]uint8

	// Tn - filetype. The characters used for these are 7-bit ASCII.
	FileType [3]uint8

	// EX = Extent counter.
	// If a file grows above 16k, then it will have multiple directory entries.
	// The first entry has EX=0, the second has EX=1 etc.
	// EX ranges from 0 to 31, thus allowing files up to 512k. CP/M 1.4 only allows 256k discs anyway.
	ExtentCounter uint8

	// S1 - reserved, set to 0.
	S1 uint8
	// S2 - reserved, set to 0.
	S2 uint8

	// RC - Number of records used in this extent.
	// 1 record is 128 bytes in length
	// If it is 80h, this extent is full and there may be another one on the disc.
	// File lengths are only saved to the nearest 128 bytes.
	RecordCount uint8

	// AL - Allocation.
	// Each AL is the number of a 1k block on the disc. If an AL number is zero,
	// that section of the file has no storage allocated to it (ie it does not exist).
	// For example, a 3k file might have allocation 5,6,8,0,0....
	// - the first 1k is in block 5, the second in block 6, the third in block 8.
	Allocation [16]uint8
}
