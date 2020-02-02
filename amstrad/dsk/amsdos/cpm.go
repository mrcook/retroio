// CP/M disk format.
//
// Amstrad computers use formats based on the standard CP/M 2 or CP/M 3
// specification, but with extensions.
package amsdos

const CpmRecordSize = 128 // CP/M records are 128 bytes in length

// DiskParameterBlock - a based on the CP/M v3 disc format, with extensions for
// the Amstrad, and referenced as the XDPB - Extended Disc Parameter Block.
//
// In general, each disk drive has an associated (16-byte) disk parameter
// header which both contains information about the disk drive and provides
// a scratchpad area for certain BIOS operations.
type DiskParameterBlock struct {
	// CP/M v3.1 DiskParameterBlock parameters

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
	//
	// Set EXM to zero if you want media compatibility with an extended
	// CP/M 1.4 system. This only applies to double-density CP/M 1.4 systems,
	// with disk sizes greater than 256K bytes.
	ExtentMask uint8

	// DSM - Storage Capacity (in blocks)
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
	// The size of the directory check vector (DRM/4)+1.
	//
	// Set bit 15 of CKS to 1 if the drive is permanently mounted.
	// Set CKS to 8000H to indicate that the drive is permanently mounted and
	// directory check-summing is not required.
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

	// PSH and PHM determine the physical sector size of the disk. All disk I/O
	// is in terms of the physical sector size. Set PSH and PSM to zero if the
	// BIOS is blocking and de-blocking instead of the BDOS.
	//
	// PSH specifies the physical record shift factor, ranging from 0 to 5,
	// corresponding to physical record sizes of 128, 256, 512, 1K, 2K, or 4K
	// bytes. It is equal to the logarithm base two of the physical record size
	// divided by 128, or LOG2(sector-size/128). See table below for PSH values.
	//
	// PHM specifies the physical record mask, ranging from 0 to 31,
	// corresponding to physical record sizes of 128, 256, 512, 1K, 2K, or 4K
	// bytes. It is equal to one less than the sector size divided by 128, or,
	// (sector-size/128)-1. See Table 3-7 for PHM values.
	//
	// Sector size | PSH | PHM
	//     128     |  0  |  0
	//     256     |  1  |  1
	//     512     |  2  |  3
	//    1024     |  3  |  7
	//    2048     |  4  | 15
	//    4096     |  5  | 31
	PhysicalShift uint8 // PSH ; Physical sector shift
	PhysicalMask  uint8 // PHM ; Physical sector mask

	//
	// Amstrad extended parameters
	//

	// Type of disc media (sidedness)
	//
	// Bit | Description
	// 0-1   0 => Single sided
	//       1 => Double sided, flip sides
	//          ie track   0 is cylinder   0 head 0
	//             track   1 is cylinder   0 head 1
	//             track   2 is cylinder   1 head 0
	//             ...
	//             track n-1 is cylinder n/2 head 0
	//             track   n is cylinder n/2 head 1
	//       2 => Double sided, up and over
	//          ie track   0 is cylinder 0 head 0
	//             track   1 is cylinder 1 head 0
	//             track   2 is cylinder 2 head 0
	//             ...
	//             track n-2 is cylinder 2 head 1
	//             track n-1 is cylinder 1 head 1
	//             track   n is cylinder 0 head 1
	//  6    Set if the format is for a high-density disc
	//         This is an extension in PCW16 CP/M, BIOS 0.09+.
	//         It is not an official part of the spec.
	//  7    Set if the format is double track.
	MediaType uint8

	// tracks/side
	TrackCountPerSide uint8

	// sectors/track
	SectorCountPerTrack uint8

	// first physical sector number
	FirstSectorNumber uint8

	// sector size, bytes
	SectorSize uint16

	// uPD765A read/write gap
	ReadWriteGap uint8

	// uPD765A format gap
	FormatGap uint8

	// MFM/Multitrack flags byte
	// Bit 7 set => Multitrack else Single track
	//     6 set => MFM mode else FM mode
	//     5 set => Skip deleted data address mark
	MultiTrackFlags uint8

	// freeze flag
	// Set to non-zero value to force this format to be used - otherwise,
	// attempt to determine format when a disc is logged in.
	FreezeFlag uint8
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

// Map for SectorSize to PSH / PHM values
//
// PSH specifies the physical record shift factor: LOG2(sector-size/128).
// PHM specifies the physical record mask:         (sector-size/128)-1.
type physicalRecords map[uint16]struct {
	PSH uint8
	PHM uint8
}

var PhysicalShiftMaskTable = physicalRecords{
	128:  {PSH: 0, PHM: 0},
	256:  {PSH: 1, PHM: 1},
	512:  {PSH: 2, PHM: 3},
	1024: {PSH: 3, PHM: 7},
	2048: {PSH: 4, PHM: 15},
	4096: {PSH: 5, PHM: 31},
}

// Disk Directory
//
// The directory is the same header as with the CP/M 2.2 standard.
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
