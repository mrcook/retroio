// CP/M 4.1 Disc Format
//
// Version 4.1 (DOS Plus) allows the use of two filesystems - CP/M and DOS.
// Those supplied with the Amstrad PC1512 cannot handle larger floppies than
// 360k (CP/M) / 1.2Mb (DOS), or larger hard drive partitions than 32Mb.
//
// For implementation details (ie: how CP/M calls behave on a DOS system) refer
// to the DOS filesystem documentation.
//
// # DOS filesystem
// The DOS filesystem can be either FAT12 or FAT16.
// The format is exactly as in PCDOS 2.11, except:
//
// Byte 0Ch of the directory entry (which is not used by PCDOS) holds the four
// "user attributes" F1'-F4':
//   Bit 7: F1'
//       6: F2'
//       5: F3'
//       4: F4'
// DRDOS-style passwords are not supported.
//
// Reference: http://www.seasip.info/Cpm/format41.html
package cpm4

import (
	"retroio/cpm/cpm3"
)

// Disk Parameter Header
//
// The DPH retains the CPM 3.1 specification.
type DiskParameterHeader struct {
	cpm3.DiskParameterHeader
}

// Disk Parameter Block
//
// The DPB retains the CPM 3.1 specification.
type DiskParameterBlock struct {
	cpm3.DiskParameterBlock // CP/M 3.1 standard
}

// Disk Directory
//
// The CP/M 4.1 directory is the same header as with the 3.1 standard,
// except for minor changes to the usage of the S1 and RC bytes.
type Directory struct {
	// NOTE the following changes against the 3.1 standard:
	//
	// BC - Last Record Byte Count
	//
	// RC - if RC is 80h, this extent is full and there may be another one on the disc.
	// File lengths are optionally saved exactly (using the BC byte).
	// This is mainly apparent in files copied from DOS-formatted media.
	cpm3.Directory
}

// Disc Label
//
// The CP/M 4.1 disc label is the same header as with the 3.1 standard,
// except for the use of the media type byte.
type DiscLabel struct {
	// NOTE the following changes against the 3.1 standard:
	//
	// MT - Media type.
	// If you search for the disc label on a DOS-formatted drive,
	// a fake disc label will be generated with MT set to 80h.
	cpm3.DiscLabel
}

// Password Control
//
// The CP/M 4.1 disc label is the same header as with the 3.1 standard.
type PasswordControl struct {
	cpm3.PasswordControl
}
