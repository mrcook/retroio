// CP/M 3.1 Disc Format
//
// Version 3.1 uses a very similar system to CP/M 2.2, but with even more formats supported.
//
// The reserved tracks will contain an image of CP/M 2.2, used when the system
// is rebooted. Discs can be formatted as data only discs, in which case they
// have no system tracks and cannot be used to reboot the system.
//
// Reference: http://www.seasip.info/Cpm/format31.html
package cpm3

import (
	"retroio/cpm/cpm1"
)

// Disk Parameter Header (Table Format) - used by the BIOS/BDOS
//
// In general, each disk drive has an associated (16-byte) disk parameter
// header which both contains information about the disk drive and provides
// a scratchpad area for certain BIOS operations.
type DiskParameterHeader struct {
	// XLT - Logical Translation Vector Offset
	//
	// Set the XLT field to the address of the logical to physical sector
	// translation table. If there is no sector translation and the logical and
	// physical sector numbers are the same, set XLT to 0000H. Disk drives with
	// identical sector skew factors can share the same translation table.
	//
	// XLT is the value passed to SECTRN in registers DE. Usually the
	// translation table consists of one byte per physical sector. Generally,
	// it is advisable to keep the number of physical sectors per logical track
	// to a reasonable value to prevent the translation table from becoming too
	// large. In the case of disks with multiple heads, you can compute the
	// head number from the track address rather than the sector address.
	TranslationVectorOffset uint16

	// These 72 bits (9 bytes) of zeroes are the scratch area the BDOS uses to
	// maintain various parameters associated with the drive.
	ScratchPad [9]uint8

	// MF is the Media Flag.
	//
	// The BDOS resets MF to zero when the drive is logged in. The BIOS can set
	// this flag and @MEDIA in the SCB to 0FFH if it detects that a drive door
	// has been opened. If the flag is set to 0FFH, the BDOS checks for a media
	// change prior to performing the next BDOS file operation on that drive.
	// If the BDOS determines that the drive contains a new volume, the BDOS
	// performs a login on that drive, and resets the MF flag to 00H. Note that
	// the BDOS checks this flag only when a system call is made, and not
	// during an operation. Usually, this flag is used only by systems that
	// support door-open interrupts.
	MediaFlag uint8

	// DPB - Disk Parameter Block Offset
	//
	// Set the DPB field to the address of a Disk Parameter Block that
	// describes the characteristics of the disk drive. Several Disk Parameter
	// Headers can address the same Disk Parameter Block if their drive
	// characteristics are identical.
	DpbOffset uint16

	// CSV - Change disk scratchpad offset
	//
	// CSV is the address of a scratchpad area used to detect changed disks.
	// This address must be different for each removable media Disk Parameter
	// Header. There must be one byte for every 4 directory entries (or 128
	// bytes of directory). In other words, length(CSV) = (DRM/4)+l. If the
	// drive is permanently mounted, set the CKS variable in the DPB to 8000H
	// and set CSV to 0000H. This way, no storage is reserved for a checksum
	// vector. The checksum vector may be located in common memory or in Bank 0.
	// Set CSV to 0FFFEH for GENCPM to set up the checksum vector.
	//
	// If CKS = (DRM+1)/4, then you must reserve this number of bytes for the
	// directory check use. If CKS = 0, then no storage is required.
	ChangeDiskOffset uint16

	// ALV - Allocation information offset
	//
	// ALV is the address of the scratchpad area called the allocation vector,
	// which the BDOS uses to keep disk storage allocation information. This
	// area must be unique for each drive.
	//
	// The allocation vector usually requires 2 bits for each block on the
	// drive. Thus, length(ALV) = (DSM/4) + 2. (See Table 3-3 for an
	// explanation of the DSM field.) In the non-banked version of CP/M 3, you
	// can optionally specify that GENCPM reserve only one bit in the
	// allocation vector per block on the drive.
	// In this case, length(ALV) = (DSM/8) + 1.
	AllocationOffset uint16

	// DIRBCB
	//
	// Set DIRBCB to the address of a single directory Buffer Control Block (BCB)
	// in an unbanked system. Set DIRBCB to the address of a BCB list head in a
	// banked system.
	//
	// Set DIRBCB to OFFFEH for GENCPM to set up the DIRBCB field. The BDOS uses
	// directory buffers for all accesses of the disk directory. Several DPHs
	// can refer to the same directory BCB or BCB list head; or, each DPH can
	// reference an independent BCB or BCB list head. Section 3.3.4 describes
	// the format of the Buffer Control Block.
	DirBufferControlBlock uint16

	// DTABCB
	//
	// Set DTABCB to the address of a single data BCB in an unbanked system.
	// Set DTABCB to the address of a data BCB list head in a banked system.
	//
	// Set DTABCB to OFFFEH for GENCPM to set up the DTABCB field.
	// The BDOS uses data buffers to hold physical sectors so that it can block
	// and deblock logical 128-byte records. If the physical record size of the
	// media associated with a DPH is 128 bytes, you can set the DTABCB field of
	// the DPH to 0FFFFH, because in this case, the BDOS does not use a data
	// buffer.
	DataBufferControlBlock uint16

	// HASH
	//
	// HASH contains the address of the optional directory hashing table
	// associated with a DPH. Set HASH to 0FFFFH to disable directory hashing.
	// Set RASH to OFFFEH to make directory hashing on the drive a GENCPM
	// option. Each DPH using hashing must reference a unique hash table. If a
	// hash table is supplied, it must be 4*(DRM+1) bytes long, where DRM is
	// one less than the length of the directory. In other words, the hash
	// table must contain four bytes for each directory entry of the disk.
	Hash uint16

	// HBANK
	//
	// Set HBANK to the bank number of the hash table. HBANK is not used in
	// unbanked systems and should be set to zero. The hash tables can be
	// contained in the system bank, common memory, or any alternate bank
	// except Bank 1, because hash tables cannot be located in the Transient
	// Program Area. GENCPM automatically sets HBANK when HASH is set to 0FFFEH.
	HashBank uint8
}

// Disk Parameter Block
//
// The disc statistics are stored in a parameter block (the DPB)
//
// NOTE: the DPB is the same as with the CP/M 1.4 standard, but for the
// addition of two extra parameters.
type DiskParameterBlock struct {
	// NOTE the following changes against the 1.4. standard:
	//
	// CKS - Checksum
	// The size of the directory check vector (DRM/4)+1.
	// Set bit 15 of CKS to 1 if the drive is permanently mounted.
	// Set CKS to 8000H to indicate that the drive is permanently mounted and
	// directory check-summing is not required.
	//
	// EXM - Extent Mask
	// Set EXM to zero if you want media compatibility with an extended
	// CP/M 1.4 system. This only applies to double-density CP/M 1.4 systems,
	// with disk sizes greater than 256K bytes.
	cpm1.DiskParameterBlock

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

// Buffer Control Block
//
// A Buffer Control Block (BCB) locates physical record buffers for the BDOS.
// The BDOS uses the BCB to manage the physical record buffers during
// processing. More than one Disk Parameter Header can specify the same BCB.
// The GENCPM utility can create the Buffer Control Block.
//
// Note that the BANK and LINK fields of the Buffer Control Block are present
// only in the banked system. Therefore, the Buffer Control Block is twelve
// bytes long in the nonbanked system, and fifteen bytes long in the banked
// system. Note also that only the DRV, BUFFAD, BANK, and LINK fields need to
// contain initial values. In Figure 3-5, which shows the form of the Buffer
// Control Block, b refers to bits.
type BufferControlBlock struct {
	// DRV
	// Identifies the disk drive associated with the record contained in the
	// buffer located at address BUFFAD. If you do not use GENCPM to allocate
	// buffers, you must set the DRV field to 0FFH.
	Drive uint8

	// REC#
	// Identifies the record position of the current contents of the buffer
	// located at address BUFFAD. REC# consists of the absolute sector number
	// of the record where the first record of the directory is zero.
	RecordPosition [3]uint8

	// WFLG
	// Set by the BDOS to OFFH to indicate that the buffer contains new data
	// that has not yet been written to disk. When the data is written, the BDOS
	// sets the WFLG to zero to indicate the buffer is no longer dirty.
	// 00 Scratch byte used by BDOS.
	WriteFlag uint8

	// TRACK
	// Contains the physical track location of the contents of the buffer.
	TrackLocation uint16

	// SECTOR
	// Contains the physical sector location of the contents of the buffer.
	SectorLocation uint16

	// BUFFAD
	// Specifies the address of the buffer associated with this BCB.
	BufferAddress uint16

	// BANK
	// Contains the bank number of the buffer associated with this BCB.
	// This field is only present in banked systems.
	BankNumber uint8

	// LINK
	// Contains the address of the next BCB in a linked list, or zero if this
	// is the last BCB in the linked list. The LINK field is present only in
	// banked systems.
	LinkedListAddress uint16
}

// Disk Directory
//
// The directory is a sequence of directory entries (also called extents),
// which contain 32 bytes of the following structure:
//
// 0U F1 F2 F3 F4 F5 F6 F7 F8 T1 T2 T3 XL BC XH RC   .FILENAMETYP....
// AL AL AL AL AL AL AL AL AL AL AL AL AL AL AL AL   ................
type Directory struct {
	// 0U - User Number
	// Allows multiple files of the same name to coexist on the disc.
	// - 0-15:  Used for file, status is the user number
	// - 16-31: used for file, status is the user number (P2DOS) or used for password extent (CP/M 3 or higher)
	// - 32:    Disc label
	// - 33:    time stamp (P2DOS)
	// - 0xE5:  File deleted
	UserNumber uint8

	// Fn - filename
	// Tn - file type extension name
	//
	// May consist of any printable 7 bit ASCII character but: < > . , ; : = ? * [ ].
	// The file name must not be empty, the extension may be empty.
	// Both are padded with blanks.
	//
	// The highest bit of each character of the file name and extension are
	// used as attributes, and have the following meanings:
	// - F0: requires set wheel byte (Backgrounder II)
	// - F1: public file (P2DOS, ZSDOS), forground-only command (Backgrounder II)
	// - F2: date stamp (ZSDOS), background-only commands (Backgrounder II)
	// - F7: wheel protect (ZSDOS)
	// - T1: File is read-only.
	// - T2: Is a system file ("hidden" on other systems). System files with User Number 0 can be read from any user number.
	// - T3: File has been backed up.
	//
	// The wheel byte is (by default) the memory location at 0x4b. If it is zero,
	// only non-privileged commands may be executed.
	Filename [8]uint8
	FileType [3]uint8

	// XL and XH store the extent number.
	//
	// A file may use more than one directory entry, if it contains more blocks
	// than an extent can hold. In this case, more extents are allocated and
	// each of them is numbered sequentially with an extent number. If a
	// physical extent stores more than 16k, it is considered to contain
	// multiple logical extents, each pointing to 16k data, and the extent
	// number of the last used logical extent is stored.
	//
	// Note: Some formats decided to always store only one logical extent in a
	// physical extent, thus wasting extent space.
	//
	// CP/M 2.2 allows 512 extents per file, CP/M 3 and higher allow up to 2048.
	//
	// Bit 5-7 of XL are 0, bit 0-4 store the lower bits of the extent number.
	// Bit 6 and 7 of XH are 0, bit 0-5 store the higher bits of the extent number.
	//
	// RC and BC determine the length of the data used by this extent.
	// The physical extent is divided into logical extents, each of them being
	// 16k in size (a physical extent must hold at least one logical extent,
	// e.g. a blocksize of 1024 byte with two-byte block pointers is not allowed).
	// RC stores the number of 128 byte records of the last used logical extent.
	// BC stores the number of bytes in the last used record.
	// The value 0 means 128 for backward compatibility with CP/M 2.2, which did
	// not support BC.
	ExtentLow   uint8 // XL as CP/M 2.2
	ByteCount   uint8 // BC - ONLY CP/M 4.1
	ExtentHigh  uint8 // XH as CP/M 2.2
	RecordCount uint8 // RC as CP/M 2.2

	// AL - Allocation stores block pointers.
	// Each AL is the number of a block on the disc. If an AL number is zero,
	// that section of the file has no storage allocated to it (ie it does not
	// exist). For example, a 3k file might have allocation 5,6,8,0,0.... - the
	// first 1k is in block 5, the second in block 6, the third in block 8.
	// AL numbers can either be 8-bit (if there are fewer than 256 blocks on the
	// disc) or 16-bit (stored low byte first).
	//
	// If the disk capacity is less than 256 blocks, AL is interpreted as 16
	// byte-values, otherwise as 8 double-byte-values. A block pointer of 0
	// marks a hole in the file. If a hole covers the range of a full extent,
	// the extent will not be allocated. In particular, the first extent of a
	// file does not necessarily have extent number 0.
	// A file may not share blocks with other files, as its blocks would be
	// freed if the other files is erased without a following disk system reset.
	// CP/M returns EOF when it reaches a hole, whereas UNIX returns zero-value
	// bytes, which makes holes invisible.
	Allocation [16]uint8
}

// Disc Label
//
// CP/M Plus support disc labels, which are stored in an arbitrary directory entry.
//
// 20 F1 F2 F3 F4 F5 F6 F7 F8 T1 T2 T3 LB PB RR RR     LABENAMETYP....
// P1 P2 P3 P4 P5 P6 P7 P8 D1 D1 D1 D1 D2 D2 D2 D2    ................
type DiscLabel struct {
	// 0x20 - Status/characteristic number of a disc label
	CharacteristicNumber uint8

	// F1-F8, T1-T3 - Label name, 7-bit ASCII
	LabelName [8]uint8

	// LB - Label mode. Bit 0 set => Label exists
	//                  Bit 4 set => Time stamp on create
	//                  Bit 5 set => Time stamp on update
	//                  Bit 6 set => Time stamp on access
	//                  Bit 7 set => Password protection enabled
	// NOTE: bits 4 and 6 are mutually exclusive!
	LabelMode uint8

	// PB - Used to decode the label password
	// To decode the password, xor this byte with the password bytes in reverse order.
	// To encode a password, add its characters to get the decode byte.
	PasswordDecode uint8

	// RR - Reserved, set to zero.
	Reserved uint8

	// MT - Media type. NOTE: only used for CP/M 4.1.
	// If you search for the disc label on a DOS-formatted drive,
	// a fake disc label will be generated with MT set to 80h.
	MediaType uint8

	// P1-P8 - password, rather feebly encrypted.
	Password [8]uint8

	// D1 - Label create date stamp
	CreatedDate Timestamp

	// D2 - Label update date stamp
	UpdatedDate Timestamp
}

// Native Time stamps
//
// P2DOS and CP/M Plus support time stamps, which are stored in each fourth directory entry.
// This entry contains the time stamps for the extents using the previous three directory entries.
// Note that you really have time stamps for each extent, no matter if it is the first extent of a
// file or not. The structure of time stamp entries is:
//
// 21 D1 D1 D1 D1 D2 D2 D2 D2 M1 00 D3 D3 D3 D3 D4    !...............
// D4 D4 D4 M2 00 D5 D5 D5 D5 D6 D6 D6 D6 M3 00 00    ................
//
// A time stamp consists of two dates: Creation and modification date (the latter being recorded when the file is closed).
// CP/M Plus further allows optionally to record the access instead of creation date as first time stamp.
type Timestamps struct {
	// 0x21 - Status/characteristic number of a date stamp.
	CharacteristicNumber uint8

	// D1 - File 1 create OR access date
	// D2 - File 1 update date
	File1 [2]Timestamp

	// M1 - File 1 password mode
	PasswordModeFile1 uint16

	// D3 - File 2 create OR access date
	// D4 - File 2 update date
	File2 [2]Timestamp

	// M2 - File 2 password mode
	PasswordModeFile2 uint16

	// D5 - File 3 create OR access date
	// D6 - File 3 update date
	File3 [2]Timestamp

	// M3 - File 3 password mode
	PasswordModeFile3 uint16

	// 00 - Reserved.
	Reserved uint8
}

// Date and Time stamp
type Timestamp struct {
	// Julian day number (little-endian).
	// Day 1 = 1978-01-01
	Day uint16

	// BCD hour, eg 13h => 13:xx
	Hour uint8

	// BCD minute
	Minute uint8
}

// Password Control
//
// CP/M Plus supports passwords, which are stored in an arbitrary directory entry.
//
// 1U F1 F2 F3 F4 F5 F6 F7 F8 T1 T2 T3 PM PB RR RR   .FILENAMETYP....
// P1 P2 P3 P4 P5 P6 P7 P8 RR RR RR RR RR RR RR RR   ................
//
// Password encryption system - this system is extremely simple:
//
// When making the password, add all 8 bytes together (packing with spaces if
// necessary). This becomes PB (the decode byte). XOR each byte with PB and
// store them backwards in the directory (ie the last byte becomes P1).
//
// To decode the password, XOR PB with the 8 bytes of the password and read it
// off backwards.
type PasswordControl struct {
	// 1U = User number + 16 (ie 16-31).
	// The user number will be the number of the file to which the password belongs.
	UserNumber uint8

	// F1-F8 - Filename of the file to which the password belongs
	// T1-T3 - File type extension of the file to which the password belongs
	Filename      [8]uint8
	FileExtension [3]uint8

	// PM - Password mode byte
	//
	// Bit: 7 set => Password required to read from file
	//      6 set => Password required to write to file
	//      5 set => Password required to delete file
	PasswordMode uint8

	// PB - Used to decode the password
	PasswordDecode uint8

	// RR - Reserved, set to 0.
	Reserved1 [2]uint8

	// P1-P8 - The password, rather feebly encrypted.
	Password [8]uint8

	// RR - Reserved, set to 0.
	Reserved2 [8]uint8
}
