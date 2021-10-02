package amsdos

// TODO: why consts? why not just set the default value on initialisation?
// Default DPB values for the Amstrad CPC SSSD disk format.
const (
	ExtentMask      uint8  = 0
	TrackCount      uint8  = 40
	SectorsPerTrack uint8  = 9
	BytesPerSector  uint16 = 512
	ReadWriteGap    uint8  = 0x2A // uPD765A read/write gap
	FormatGap       uint8  = 0x52 // uPD765A format gap

	DiskByteSize = int(TrackCount) * int(SectorsPerTrack) * int(BytesPerSector)

	BLS uint16 = 1024
	DSM uint16 = uint16(DiskByteSize / int(BLS))
	DRM uint16 = 64 // All formats have 64 directory entries: 64 * 32-bytes = 2048 bytes.
)

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

	// MFM/Multi-track flags byte
	// Bit 7 set => Multi-track else Single track
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
