# CP/M disk and file system format.

Sources:
- https://www.cpm8680.com/cpmtools/cpm.htm
- http://www.cpm.z80.de/randyfiles/DRI/CPM-86_System_Guide.pdf
- http://www.sydneysmith.com/wordpress/1651/cpm-disk-formats/
- http://www.stjarnhimlen.se/apple2/CPM.ref.txt
- Wikipedia


## CP/M Disc Layout

* Zero or more reserved tracks;
* One or more data blocks
  - each 1k in size
  - data blocks can span tracks and usually contain multiple sectors
* Any spare sectors - ignored by CP/M


## Drive Characteristics

    65536: 128 Byte Record Capacity
     8192: Kilobyte Drive Capacity
      128: 32 Byte Directory Entries
        0: Checked Directory Entries
     1024: Records/Extent
      128: Records/Block
       58: Sectors/Track
        2: Reserved Tracks


## Characteristics sizes

Each CP/M disk format is described by the following specific sizes:

    Sector size in bytes
    Number of tracks
    Number of sectors
    Block size
    Number of directory entries
    Logical sector skew
    Number of reserved system tracks (optional)
    Offset to start of volume (optional)

A block is the smallest allocatable storage unit.

CP/M supports block sizes of 1024, 2048, 4096, 8192 and 16384 bytes.
Unfortunately, this format specification is not stored on the disk and
there are lots of formats. Accessing a block is performed by accessing
its sectors, which are stored with the given software skew.

Device areas - A CP/M disk contains these areas:

    Volume Offset (optional)
    System tracks (optional)
    Directory
    Data

The system tracks store the boot loader and CP/M itself. In order to save
disk space, there are non-bootable formats which omit those system tracks.
The term disk capacity always excludes the space for system tracks.
Note that there is no bitmap or list for free blocks. When accessing a drive
for the first time, CP/M builds this bitmap in core from the directory.

A hard disk can have the additional notion of a volume offset to locate the
start of the drive image (which may or may not have system tracks associated
with it). The base unit for volume offset is byte count from the beginning
of the physical disk, but specifiers of K, M, T or S may be appended to
denote kilobytes, megabytes, tracks or sectors. If provided, a specifier
must immediately follow the numeric value with no whitespace. For
convenience upper and lower case are both accepted and only the first letter
is significant, thus 2KB, 8MB, 1000trk and 16sec are valid values. Offset
must appear subsequent to track, sector and sector length values.


### Sector Skewing

Skewing is the rearrangement of sectors on a disc, so that by the time the
computer has read and processed one sector, the next will be in the right
position for the disc controller to read. Otherwise the poor controller has to
wait for the disc to make a full revolution before the right sector appears.

There are two types of skewing - software and hardware.

### Hardware skewing

In this system, when the disc is formatted the sectors are laid out in the
required order - for example 1, 4, 7, 2, 5, 8, 3, 6, 9. The retrieval of the
correct sector is handled transparently by the disc controller and no special
code is needed to handle such discs.

### Software skewing

In this system, the disc is formatted with the sectors in numeric order - for
example 1, 2, 3, 4, 5, 6, 7, 8, 9. However, the software which reads/writes the
disc has a translation table which it uses whenever supplying sector numbers to
the controller - this might read 1, 3, 5, 7, 9, 2, 4, 6, 8 - so if a program
writes to what it thinks is the sixth sector of a track, this is translated to
sector number 2.

The problem with software skewing comes when another computer wants to read the
discs created by the first. Unless the contents of the translation table are
available, the information on the tracks will appear to be hopelessly jumbled up.

Under CP/M 1, the skewing is fixed for its 8" discs. The table reads:

  1,7,13,19,25,5,11,17,23,3,9,15,21,2,8,14,20,26,6,12,18,24,4,10,16,22

Under CP/M 2 and later, software skewing is handled by the SECTRAN system call.
SECTRAN requires the address of a translation table; see the description of the
Disk Parameter Header to find the address of this table.


## Disk Parameter Tables For Specific Disks


### Standard CP/M 8" SSSD disk

IBM 3740 standard 250k 8" soft-sectored SSSD (single-sided, single density),
floppy-disk.

    128 bytes/sector (one logical record)
    26 sectors/track (1-26), software skewed
    77 tracks (0-76) - 2 system tracks
    2 reserved tracks
    75 used tracks ==> 243.75 user KBytes/disk
    1024 bytes/block ==> 243 blocks/disk ==> DSM=242
    Directory in 2 first blocks ==> 64 directory entries ==> 241.75 KBytes data

    Storage/disk: 256256 bytes (77*26*128)
    File size: any number of sectors from zero to capacity of disk
    Extent: 1 kBytes - 8 sectors (smallest file space allocated)

The `reserved` tracks will contain an image of CP/M 1.4, used when the system is
rebooted. It can therefore be deduced that CP/M 1.4 fits in 6.5k.

**Sector skew table:**

1 byte/sector, 6 sectors standard (space between consecutive physical sectors on track):

    1, 7, 13, 19, 25, 5, 11, 17, 23, 3, 9, 15, 21,
    2, 8, 14, 20, 26, 6, 12, 18, 24, 4, 10, 16, 22

**System:**

    Track 0 & 1 (optional)
    Track 0 sector 1: boot loader
    Track 0 sectors 2-26:  CCP & BDOS
    Track 1 sectors 1-17:  CCP & BDOS
    Track 1 sectors 18-26: CBIOS

**Directory:**

    Track 2:
        16 sectors typical
        32 bytes/entry
        64 entries typical
        extents 0 and 1

**User file area:**

Remaining sectors on Track 2 and 3 to 76, extents 2 and above.

**A Standard CP/M 8" SSSD diskette allocation**

     Track #  | Sector # |  Page #  | Memory Address | CP/M Module Name
    ---------------------------------------------------------------------
        00    |    01    |          | (boot address) | Cold Start Loader
    ---------------------------------------------------------------------
        00    |    02    |    00    |     2900H+b    |       CCP
        00    |    03    |    00    |     2980H+b    |       CCP
        00    |    04    |    01    |     2A00H+b    |       CCP
        00    |    05    |    01    |     2A80H+b    |       CCP
        00    |    06    |    02    |     2B00H+b    |       CCP
        00    |    07    |    02    |     2B80H+b    |       CCP
        ..    |    ..    |    ..    |     .......    |       ...
        00    |    18    |    08    |     3100H+b    |       BDOS
        00    |    19    |    08    |     3180H+b    |       BDOS
        ..    |    ..    |    ..    |     .......    |       ....
        00    |    26    |    12    |     3500H+b    |       BDOS
        01    |    01    |    12    |     3580H+b    |       BDOS
        01    |    02    |    13    |     3600H+b    |       BDOS
        01    |    03    |    13    |     3680H+b    |       BDOS
        ..    |    ..    |    ..    |     .......    |       ....
        01    |    18    |    21    |     3E00H+b    |       BIOS
        01    |    19    |    21    |     3E80H+b    |       BIOS
        01    |    20    |    22    |     3F00H+b    |       BIOS
        01    |    21    |    22    |     3F80H+b    |       BIOS
    ---------------------------------------------------------------------
        01    |  22-26                               (not currently used)
    ---------------------------------------------------------------------
        02    |  01-08                                 Directory block 1
        02    |  09-16                                 Directory block 2
        02    |  17-26                                 Data
      03-76   |  01-26                                 Data

Track 0, Sector 1 is the optional software boot section.

**DPB (Disk Parameter Block)**

    SPT 16b     26    Sectors per track
    BSH  8b      3    Block shift factor
    BLM  8b      7    Block shift mask
    EXM  8b      0    Extent mask - null
    DSM 16b    242    Disk size - 1 (in blocks)
    DRM 16b     63    directory mask = dir entries - 1
    AL0  8b   0C0H    Dir Alloc 0
    AL1  8b      0    Dir Alloc 1
    CKS 16b     16    Directory check vector size
    OFF 16b      2    Track offset: 2 system tracks
    
    Dirbuf 128 bytes
    ALV     31 bytes
    CSV     16 bytes

Block size 1024 bytes ==> BSH=3, BLM=7

DSM = 242 blocks

Disk size:

* 243.75 KBytes excluding system tracks
* 250.25 KBytes including system tracks


### 5.25" CP/M Disks SSSD disk

No standard 5.25" CP/M disk format exists - there were perhaps "two dozen"
different ones on the market - although the Xerox 820 format was one of the
more popular.

By 1982 the most popular disks were 90k soft-sector SSSD's (single-sided
single density) with the following geometry:

    40 tracks
    18 sectors per track
    128 byte sectors.
    81K data storage


### CDOS 5.25" DSDD Disks

CDOS was the Cromemco Disk Operating System. It was a CP/M look alike (and behave alike).

Single-sided, single-density (SSSD) disk in the 5.25" format.

    40 tracks

    Track #0:
      * 18 sectors per track
      * 128 byte sectors

    Tracks #01 - #39:
      * 10 sectors per track
      * 512 byte sectors.

    81K storage capacity


## Apple CP/M 5.25" disks

Physical format:

      A            B                C
    
    ---- Standard -----     ----- Special ------
    13-sect     16-sect     80-trk/16-sec/2-side

    Bytes/sector          256         256               256
    Sectors/track          13          16                16
    Tracks                 35          35                80
    Heads                   1           1                 2

Sector skew table (1 byte/sector):  _no sector skew in CP/M BIOS_

    13-sector disks: hard sector skew
    16-sector disks: soft sector skew in 6502 code (CP/M RWTS)

DPB (Disk Parameter Block)

                    A       B       C
    
    SPT 16b         26      32      32      Sectors per track
    BSH  8b          3       3       4      Block shift factor
    BLM  8b          7       7      15      Block shift mask
    EXM  8b          0       0       0      Extent mask
    DSM 16b        103     127     313      Disk size - 1 (in blocks)
    DRM 16b         47      63     255      Directory mask = dir entries - 1
    AL0  8b       0C0H    0C0H    0F0H      Dir Alloc 0
    AL1  8b          0       0       0      Dir Alloc 1
    CKS 16b         12      16      64      Directory check vector size
    OFF 16b          3       3       3      Track offset: 3 system tracks
    
    Block size    1024    1024    2048
    Dir entries     48      64     256
    Dir blocks       2       2       4
    DSM+1          104     128     314 blocks
    Disk size      104     128     628 KBytes (excluding system tracks)
                   113.75  140     640 KBytes (including system tracks)
    
    Dirbuf         128     128     128 bytes
    ALV             14      17      40 bytes
    CSV             12      16      64 bytes




## CP/M 3 system organisation

    Track M:       CP/M Data Region
                   CP/M Directory Region
    Track N:       CCP (optional)
    System Tracks: CPMDR
                   Cold Boot Loader

The first N tracks are the system tracks; the remaining tracks, the data tracks,
are used by CP/M 3 for file storage. Note that the system tracks are used by
CP/M 3 only during system cold start and warm start. All other CP/M 3 disk
access is directed to the data tracks of the disk. To maintain compatibility
with Digital Research products, you should use an eight-inch, single-density,
IBM 3740 formatted disk with two system tracks.

