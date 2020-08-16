# Amstrad Disk Notes

Sources:

  - https://www.cpcwiki.eu/imgs/b/bc/S968se09.pdf
  - http://www.seasip.info/Cpm/formats.html -> http://www.seasip.info/Cpm/format22.html
  - http://www.moria.de/~michael/cpmtools/
  - http://www.seasip.info/Cpm/amsform.html

Access to the file system itself references the AMSDOS and CP/M specifications, and other tools.


## Amstrad DSK Geometry

Geometry for DSK files follows the specification at http://cpctech.cpc-live.com/docs/dsk.html

  * Track 0 (or Track 0 side 0 for double-sided disks), if track data exists,
    will immediately follow the Disc Information Block and will start at
    offset `&100` in the disc image file.
  * All tracks must have a "Track Information Block"
  * Track lengths are stored in the same order as the tracks in the image. In
    the case of a double-sided disk:
      Track 0 side 0, Track 0 side 1, Track 1 side 0, etc.
  * The track blocks are stored in increasing order `0..number of tracks`, with
    alternating sides interleaved if the disc image describes a double-sided
    disk. E.g. if the disk image represents a double-sided disk, the track order is:
    - track 0 side 0,
    - track 0 side 1,
    - track 1 side 0,
    - track 1 side 1....
    - track (number of tracks-1) side 0, track (number of tracks-1) side 1

The tracks are always ordered in this way regardless of the disc-format
described by the disc image.

A standard disk image can be used to describe a copy-protected disk, but will
often result in a file which is larger than the same disk described by an
extended disk image.

For a standard disk image to represent a copy-protected disk:

  - All track sizes in the standard disk image must be the same. This value
    therefore would be the size of the largest track, and other tracks would
    have unused space in them.
  - All sector sizes within each track must be the same size, but not
    necessarily the same size as the sectors for another track. If a track
    contained different sized sectors, the size of the largest sector should
    be used. This would result in some wasted space.


### Disc media type (sided-ness)

     Bit |  Description
    -----|------------------------------------------------
     0-1 | 0 => Single sided
         | 1 => Double sided, flip sides
         |    ie track   0 is cylinder   0 head 0
         |       track   1 is cylinder   0 head 1
         |       track   2 is cylinder   1 head 0
         |       ...
         |       track n-1 is cylinder n/2 head 0
         |       track   n is cylinder n/2 head 1
         | 2 => Double sided, up and over
         |    ie track   0 is cylinder 0 head 0
         |       track   1 is cylinder 1 head 0
         |       track   2 is cylinder 2 head 0
         |       ...
         |       track n-2 is cylinder 2 head 1
         |       track n-1 is cylinder 1 head 1
         |       track   n is cylinder 0 head 1
      6  | Set if the format is for a high-density disc
         |   This is an extension in PCW16 CP/M, BIOS 0.09+.
         |   It is not an official part of the spec.
      7  | Set if the format is double track.


### General DSK Format Geometry

**Single sided DSK images:**

  * Disc Information Block
  * Track 0 data
    - Track Information Block
    - Sector data
  * Track 1 data
    - Track Information Block
    - Sector data
  * . . . .
  * Track (number_of_tracks-1) data
    - Track Information Block
    - Sector data

**Double sided DSK images:**

  * Disc Information Block
  * Track 0 side 0 data
    - Track Information Block
    - Sector data
  * Track 0 side 1 data
    - Track Information Block
    - Sector data
  * . . . .
  * Track (number_of_tracks-1) side 1 data
    - Track Information Block
    - Sector data


## Amstrad Disc Formats

Amstrad computers use standard CP/M 2 or CP/M 3 formats. The disc formats used
also include automatic format detection systems.

AMSDOS and the CP/M 2.2 BIOS support three different disc formats: `SYSTEM`,
`DATA ONLY`, and `IBM` formats.

The CP/M Plus BIOS supports only the `SYSTEM` and `DATA` formats.

The BIOS automatically detects the format of a disc. Under CP/M this occurs
for drive `A` at a warm boot, and for drive `B` the first time it is accessed.
Under AMSDOS this occurs each time a disc with no open files is accessed. To
permit this automatic detection each format has unique sector numbers as
detailed below.

3-inch discs are double-sided, but only one side may be accessed at a time
depending on which way round the user inserts the disc. There might be
different formats on the two sides.

### Common To All Formats

  * Single sided (the two sides of a 3-inch disc are treated separately).
  * 512 byte physical sector size.
  * 40 track numbered 0 to 39.
  * 1024 byte CP/M block size.
  * 64 directory entries.

### SYSTEM Format

  * 9 sectors per track numbered #41 to #49.
  * 2 reserved tracks.
  * 2 reserved tracks.
  * 2 to 1 sector interleave.

The system format is the main format supported, CP/M can only be loaded (Cold Boot)
from a system format disc. CP/M 2.2 also requires a system format disc to warm boot.

The reserved tracks are used as follows:

    Track 0 sector   #41:       boot sector.
    Track 0 sector   #42:       configuration sector
    Track 0 sectors  #43..#47:  unused
    Track 0 sectors  #41..#49:  and
    Track 1 sectors  #48..#49:  CCP and BIOS

    CP/M Plus only uses Track 0 sector #41 as a boot sector.
    Track 0 sectors #42..#49, and Track 1 are unused.

NOTE: another format called `VENDOR` format is a special version of system
format which does not contain any software on the two reserved tracks. It
is intended for use in software distribution.

### DATA ONLY Format

  * 9 sectors per track numbered #C1 to #C9.
  * 0 reserved tracks.
  * 2 to 1 sector interleave.

This format is not recommended for use with CP/M 2.2 since it is not possible
to ‘warm boot’ from it. However, because there is a little more disc space
available it is useful for AMSDOS or CP/M Plus.

### IBM Format

  * 8 sectors per track numbered 1 to 8
  * 1 reserved track
  * no sector interleave

This format is logically the same as the single-sided format used by CP/M on
the IBM PC. It is intended for specialist use and is not otherwise recommended
as it is not possible to 'warm boot' from it.


## AMSDOS and non-CP/M disc formats

Amstrad and Locomotive Software have made a number of non-CP/M systems which
are based on the CP/M 2 disc format. These are:

  * `LocoScript` - word processor with built-in operating system.
  * `AMSDOS`     - disc operating system for the CPC computers.
  * `+3DOS`      - disc operating system for the Spectrum +3.

Some common Amstrad Extended Disk Parameter Block headers for the different
disk formats (all numbers are in Hex):

                  |       CPC       |       PCW       |  PCW16
      XDPB field   | System |  Data  |  180k  |  720k  |  1.4Mb
    --------------+--------+--------+--------+--------+--------
     SPT          |    24  |    24  |    24  |    24  |    48
     BSH          |    03  |    03  |    03  |    04  |    05
     BLM          |    07  |    07  |    07  |    0F  |    1F
     EXM          |    00  |    00  |    00  |    00  |    01
     DSM          |  00AA  |  00B3  |  00AE  |  0164  |  0164
     DRM          |    3F  |    3F  |    3F  |    FF  |    FF
     AL0          |    C0  |    C0  |    C0  |    F0  |    C0
     AL1          |    00  |    00  |    00  |    00  |    00
     CKS          |    10  |    10  |    10  |    40  |    40
     OFF          |    02  |    00  |    01  |    01  |    01
     PSH          |    02  |    02  |    02  |    02  |    02
     PHM          |    03  |    03  |    03  |    03  |    03
     Sidedness    |    00  |    00  |    00  |    81  |    C1
     Cylinders    |    28  |    28  |    28  |    50  |    50
     Sectors      |    0A  |    0A  |    0A  |    0A  |    12
     1st Phys Sec |    41  |    C1  |    01  |    01  |    01
     Sector Size  |  0200  |  0200  |  0200  |  0200  |  0200
     R/W Gap      |    2A  |    2A  |    2A  |    2A  |    1B
     Format Gap   |    52  |    52  |    52  |    52  |    54
     MFM Mode     |    60  |    60  |    60  |    60  |    60
     Freeze Flag  |    00  |    00  |    00  |    00  |    FF



### Disc Parameter Block

Amstrad CP/M (and +3DOS) has an eXtended Disc Parameter Block (XDPB).
The DPB is not stored on disc.

This simple system is used by CPC computers if the first physical sector is:

    41h - A System formatted disc:

          single sided, single track,
          40 tracks, 9 sectors/track, 512-byte sectors,
          2 reserved tracks, 1k blocks,
          2 directory blocks,
          gap lengths 2Ah and 52h,
          bootable

    C1h - A Data formatted disc:

          single sided, single track,
          40 tracks, 9 sectors/track, 512-byte sectors,
          no reserved tracks, 1k blocks,
          2 directory blocks,
          gap lengths 2Ah and 52h,
          not bootable

In addition to the XDPB system, the PCW and Spectrum +3 can determine the format
of a disc from a 16-byte record on track 0, head 0, physical sector 1.

If all bytes of the spec are `0E5h`, it should be assumed that the disc is a
173k PCW/Spectrum +3 disc, ie:

    single sided, single track
    40 tracks, 9 sectors/track, 512-byte sectors
    1 reserved track, 1k blocks
    2 directory blocks
    gap lengths 2Ah and 52h
    not bootable

#### PCW16 extended boot record

The "boot record" system has been extended in PCW16 CP/M (BIOS 0.09 and later).
The extension is intended to allow a CP/M _partition_ on a DOS-formatted floppy disc.

An extended boot sector (cylinder 0, head 0, sector 1) has the following characteristics:

  - First byte is `0E9h` or `0EBh`
  - Where DOS expects the disc label to be (at `sector + 2Bh`) there are 11 ASCII bytes
    of the form `CP/M????DSK`, where `?` can be any character.
  - At `sector + 7Ch` are the four ASCII bytes `CP/M`
  - At `sector + 80h` is the disc specification as described above.


### Boot Sector

In order that non-CP/M systems may be implemented at a later date the BIOS
initialization is performed, in part, by a boot program which is read from the
disc before attempting to load CP/M. In the non-CP/M case the boot program
would not jump to the warm boot routine but go on its own way, using the BIOS
and firmware routines as desired.

The boot program is in the boot sector which is the first sector (sector #41)
on track 0.


## Amstrad disk formats

Source: https://www.cpcwiki.eu/imgs/b/bc/Knife_Plus_Manual.pdf

### Amstrad IBM Format

    Single-sided
    512-byte physical sector size
    40 tracks numbered 0 to 39
    1024-byte block size
    64 directory entries
    8 sectors per track numbered #01 to #08
    1 reserved track

### Amstrad System / Vendor Format (CPC6128)

    Single-sided
    512-byte physical sector size
    40 tracks numbered 0 to 39
    1024-byte block size
    64 directory entries
    9 sectors per track numbered #41 to #49
    2 reserved tracks

### Amstrad Data Format (CPC6128)

    Single-sided
    512-byte physical sector size
    40 tracks numbered 0 to 39
    1024-byte block size
    64 directory entries
    9 sectors per track numbered #C1 to #C9
    0 reserved tracks

### Amstrad CF2 Format (PCW8256/PCW8512 drive A format)

    Single-sided
    512-byte physical sector size
    40 tracks numbered 0 to 39
    1024-byte block size
    64 directory entries
    9 sectors per track numbered #09 to #09
    2 reserved tracks*

NOTE: 'reserved tracks' in the PDF may be wrong, it might only be 1 reserved track.

### Amstrad CF2-DD Format

PCW8256/8512 drive B format, PCW9512 drive A format.

    Double-sided
    512-byte physical sector size
    160 tracks numbered 0 to 159
    2048-byte block size
    256 directory entries
    9 sectors per track numbered #09 to #09
    1 reserved track


## Extended DSK image definition

Source: http://cpctech.cpc-live.com/docs/extdsk.html

The extended DSK image is a file designed to describe copy-protected floppy disk
software. It's definition was defined by Marco Vieth, Ulrich Doewich and Kevin
Thacker.

This format has been widely adopted and is one of the major file formats used
(the other major format is the standard disk image with the `MV - CPC`
identifier).

The extended disk image format should be used for copy-protected disc software,
or in place of a standard disk image if the resulting image is smaller than in
the standard disk image form.


### Extended DiSK Format (Rev.5)

The disc image has the following format:

    DISK INFORMATION BLOCK (256 bytes)
    TRACK INFORMATION BLOCK * number of tracks * number of sides

The track blocks are stored in increasing order 0..number of tracks, with
alternating sides interleaved if the disc image describes a double-sided disk.
e.g. if the disk image represents a double-sided disk, the order of tracks is:

    track 0 side 0
    track 0 side 1
    track 1 side 0
    track 1 side 1....
    track (number of tracks-1) side 0
    track (number of tracks-1) side 1

The tracks are **always** ordered in this way regardless of the disc-format
described by the disc image.

The location of the track information block is found by using the track size
table.


### DISK INFORMATION BLOCK

The "DISK INFORMATION BLOCK" is always located at offset 0 in the disk image
file, and has the following structure:

    offset    description                               bytes
    --------+------------------------------------------+-----
    00 - 21   "EXTENDED CPC DSK File\r\nDisk-Info\r\n"   34
    22 - 2f   name of creator (utility/emulator)         14
    30        number of tracks                           1
    31        number of sides                            1
    32 - 33   unused                                     2
    34 - xx   track size table                           number of tracks*number of sides

NOTES:

  * An extended DSK image is identified by the `EXTENDED` tag.
  * The track size at offset `32h` and `33h`, used by the STANDARD disk image is
    ignored for extended format DSK images.
  * If track data exists, then it starts at offset `100h`.
  * The `EXTENDED` tag is present to prevent existing emulators which support the
    standard DSK image from interpreting the data wrong and possibly crashing.
  * `\r` is the C programming language equivalent of ASCII character `13`.
  * `\n` is the C programming language equivalent of ASCII character `10`.


### TRACK OFFSET TABLE

    offset   description                       bytes
    -------+----------------------------------+-----
    01       high byte of track 0 length        1
             (equivalent to track length/256)
    ...      ...                                ...

NOTES:

  * Depending on the information in the disk information block, the table contains
    - track lengths for a single sided floppy disc
    - track lengths for a double-sided floppy disc
  * track lengths are stored in the same order as the tracks in the image e.g. In
    the case of a double-sided disk:
      Track 0 side 0, Track 0 side 1, Track 1 side 0 etc...
  * A size of "0" indicates an unformatted track. In this case there is no data,
    and no track information block for this track in the image file!
  * Actual length of track data = (high byte of track length) * 256
  * Track length includes the size of the TRACK INFORMATION BLOCK (256 bytes)
  * The location of a Track Information Block for a chosen track is found by
    summing the sizes of all tracks up to the chosen track plus the size of
    the Disc Information Block (`&100` bytes). The first track is at offset &100
    in the disc image.


### TRACK INFORMATION BLOCK

    offset    description              bytes
    --------+-------------------------+-----
    00 - 0c   "Track-Info\r\n"          13
    0d - 0f   unused                    3
    10        track number              1
    11        side number               1
    12 - 13   unused                    2
    14        sector size               1
    15        number of sectors         1
    16        GAP#3 length              1
    17        filler byte                1
    18 - xx   Sector Information List   xx

NOTES:

  * Identical to the original DSK format definition.


### SECTOR INFORMATION LIST

    offset    description                                                      bytes
    --------+-----------------------------------------------------------------+-----
    00        track (equivalent to C parameter in NEC765 commands)              1
    01        side (equivalent to H parameter in NEC765 commands)               1
    02        sector ID (equivalent to R parameter in NEC765 commands)          1
    03        sector size (equivalent to N parameter in NEC765 commands)        1
    04        FDC status register 1 (equivalent to NEC765 ST1 status register)  1
    05        FDC status register 2 (equivalent to NEC765 ST2 status register)  1
    06 - 07   actual data length in bytes                                       2

NOTES:

  * Identical to the original definition except for the addition of the sector
    data length. This value is in bytes and stored in little endian notation.
    (low byte followed by high byte)
  * The location of each sectors data is found by adding the size of the previous
    sectors, plus the size of the 256 byte header.
  * For 8k Sectors (`N="6"`), only `1800h` bytes are stored. Please see extensions below!!!


### Extensions to the above specification

**This extension was proposed by John Elliott:**

    Extension to TRACK INFORMATION BLOCK:

    offset  description                             bytes
    ------+----------------------------------------+-----
    12      Data rate. (See note 1 and note 3)       1
    13      Recording mode. (See note 2 and note 3)  1

Notes:

1 - Data rate defines the rate at which data was written to the track. This
   value applies to the entire track.

    Date rate   description
    ----------+------------
        0       Unknown
        1       Single or double density
        2       High Density
        3       Extended density

Existing files should have zeroes in these bytes; hence the use of 0 for unknown.


2 - Recording mode is used to define how the data was written. It defines the
   encoding used to write the data to the disc and the structure of the data
   on the disc including the layout of the sectors. This value applies to the
   entire track.

    Date rate   description
    ----------+------------
        0       Unknown
        1       FM
        2       MFM

Existing files should have zeroes in these bytes; hence the use of 0 for unknown.


3 - How to determine the actual rate.

The NEC765 floppy disc controller is supplied with a single clock. When reading
from and writing to a disc using the NEC765 you can choose FM or MFM recording
modes. Use of these modes and the clock into the NEC765 define the final rate
at which the data is written to the disc.

When FM recording mode is used, data is read from or written to at a rate which
is double that of when MFM is used. The time for each bit will be twice the
time for MFM.

Examples:

    NEC765 Clock  FM/MFM  Actual rate
    4Mhz  FM  4us per bit
    4Mhz  MFM 2us per bit


**This extension was proposed by Simon Owen:**

1. It has been found that many protections using 8K Sectors (`N="6"`) do store
   more than `&1800` bytes of usable data. It was thought that `&1800` was the
   maximum usable limit, but this has proved wrong. So you should support 8K
   of data to ensure this data is read correctly. The size of the sector will
   be reported in the SECTOR INFORMATION LIST as described above.

   For sector size `N="7"` the full 16K will be stored. It is assumed that
   sector sizes are defined as 3 bits only, so that a sector size of `N="8"`
   is equivalent to `N="0"`.

2. Storing Multiple Versions of Weak/Random Sectors.

   Some copy protections have what is described as 'weak/random' data. Each
   time the sector is read one or more bytes will change, the value may be
   random between consecutive reads of the same sector.

   To support these formats the following extension has been proposed.

   Where a sector has weak/random data, there are multiple copies stored. The
   actual sector size field in the SECTOR INFORMATION LIST describes the size
   of all the copies. To determine if a sector has multiple copies then compare
   the actual sector size field to the size defined by the N parameter. For
   multiple copies the actual sector size field will have a value which is a
   multiple of the size defined by the N parameter. The emulator should then
   choose which copy of the sector it should return on each read.
