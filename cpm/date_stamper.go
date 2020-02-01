// CP/M disk and file system format.
//
// References: https://www.cpm8680.com/cpmtools/cpm.htm
package cpm

// DateStamper Time stamps
// TODO: check this spec, especially about the checksum/checksums.
//
// The DateStamper software added functions to the BDOS to manage time stamps
// by allocating a read only file with the name "!!!TIME&.DAT" in the very
// first directory entry, covering the very first data blocks.
//
// F0 F1 F2 F3 F4 F5 F6 F7 DC DC DC DC DC DA DA DA    !!!TIME.........
// DA DA DM DM DM DM DM CS DC DC DC DC DC DA DA DA    ................
// DA DA DM DM DM DM DM CS DC DC DC DC DC DA DA DA    ................
// DA DA DM DM DM DM DM CS DC DC DC DC DC DA DA DA    ................
//
// A timestamp file contains the 7-byte magic number, followed by seven 16-byte timestamp blocks, and a checksum.
type DateStamperTimestamps struct {
	// The magic number is used for the first 7 entries of each 128-byte record
	// and contains the characters !, !, !, T, I, M and E (!!!TIME).
	MagicNumber [7]uint8

	Stamps [7]DateStamperDateField

	// TODO: is this an additional checksum, or just the checksum of the last Timestamp entry?

	// The check-sum is used on every 8th entry (last entry in 128-byte record)
	// and is the sum of the first 127 bytes of the record.
	Checksum uint8
}

// DateStamper Timestamp
//
// A timestamp contains one 16-byte block per directory entry.
type DateStamperTimestamp struct {
	Create   DateStamperDateField
	Access   DateStamperDateField
	Modify   DateStamperDateField
	Checksum uint8
}

// DateStamper DateField
//
// All values are in BCD format.
type DateStamperDateField struct {
	// Year: no century, so it is sane assuming any year < 70 means 21st century
	Year uint8
	// Month
	Month uint8
	// Date
	Day uint8

	// If the high bit is set on Hour, the Hour/Minute bytes are used as the
	// High/Low bytes of the counter for systems without real time clock.
	Hour   uint8
	Minute uint8
}
