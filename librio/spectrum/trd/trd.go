// https://web.archive.org/web/20190912212218/http://www.zx-modules.de:80/fileformats/trdformat.html
// https://faqwiki.zxnet.co.uk/wiki/TR-DOS_filesystem
// http://www.8bit-wiki.de/index.php?id=3&filt=Sinclair/ZX_Spectrum/interfaces/betadisk/_manual/&cid=16062&mode=dl&tx=330ea210000 (PDF)
package trd

// The TR-DOS operating system uses the filesystem of the Beta Disc and
// Beta 128 Disc disk interfaces from Technology Research Ltd.
//
// Floppy disks formatted to 16 sectors per track, with 256-byte sectors.
// Either single-sided or double-sided disks may be used, and 40 or 80 tracks
// may be used per side.
//
// The first track of the disk (h0t0s1..h0t0s15) contains the file descriptors
// (h0t0s1..h0t0s8) and the disk info (h0t0s9). The remaining space on this
// track is unused. Space from h0t1s1 on a single sided disk or from h1t0s1 on
// a double sided disk can be allocated to files. The files are not fragmented.
//
// TRD disk image format:
//     Offset | Type          | Length | Description
//         0    DIRECTORY       2048     file allocation table
//      2048    SPECIFICATION    256     disk specification
//      2304    BYTE[]          1792     filled with zero (filler)
//      4096    BYTE[]           256     16th data sector (start of data sectors)
//      4352    BYTE[]           256     17th data sector
//      4608    BYTE[]           256     18th data sector
//      ....
//      ????    BYTE[]           256     last data sector
//
// The first track (16 sectors) is reserved and contains the FAT and the disk specification.
//
// The index of data sectors begins at sector #16. The data sectors contain the
// data of file bodies. The maximum count of sectors depends on the disk type.
type TRD struct {
	directory     Directory     // file allocation table
	specification Specification // disk specification
	unusedFiller  [1792]uint8   // filler: filled with 0x00
	files         []File        // data sectors (file bodies)
}
