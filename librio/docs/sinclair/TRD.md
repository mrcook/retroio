# TRD disk image format



Format name: TR-DOS image format, Format creator: RamSoft

A TR-DOS disk contains 1 or 2 sides with 40 or 80 tracks per side and with 16
sectors per track. Each sector is 256 bytes long.

There are 4 different disk formats, described below.

                            total     total    total    reserved   reserved
    disk type     | sides | tracks | sectors | bytes  | sectors  | bytes
    DS, 80 tracks     2      160      2560     655360       16     4096
    DS, 40 tracks     2      80       1280     327680       16     4096
    SS, 80 tracks     1      80       1280     327680       16     4096
    SS, 40 tracks     1      40        640     163840       16     4096

The counting of the sectors includes the reserved sectors, that means that
the first writeable sector has the index number 16.


## TRD file format (all 4 disk types)

    Offset | Type          | Length | Description
        0    DIRECTORY       2048     file allocation table
     2048    SPECIFICATION    256     disk specification
     2304    BYTE[]          1792     filled with zero (filler)
     4096    BYTE[]           256     16th data sector (start of data sectors)
     4352    BYTE[]           256     17th data sector
     4608    BYTE[]           256     18th data sector
     ....
     ????    BYTE[]           256     last data sector

The first track (16 sectors) is reserved and contains the FAT and the disk specification.

The index of data sectors begins at sector #16. The data sectors contain the
data of file bodies. The maximum count of sectors depends on the disk type.


## Directory (2048 bytes)

    Offset | Type          | Length | Description
        0    DIRENTRY #0       16     header of 1st file
       16    DIRENTRY #1       16     header of 2nd file
       32    DIRENTRY #2       16     header of 3rd file
     ....
     2032    DIRENTRY #127     16     header of 128th file

The directory always contains 128 entries, but the list ends when the first
byte of the directory entry (= of the file name) is #0. Therefore, any
following entries are invalid and may contain random data.


### Directory Entry (16 bytes)

    Offset | Type   | Length | Description              | Additional Information
       0     CHAR[]      8     file name                  case-sensitive; if the first character is
                                                            #0, then it's the end of the directory.
                                                            #1 indicates a deleted file, which is
                                                            still present on the disk.
       8     CHAR        1     file extension             character that describes the file type:
                                                            "B" = Basic program
                                                            "D" = DATA array (numeric or alphanumeric)
                                                            "C" = CODE 
                                                            "#" = Print file (may be split into several sub-files with max. 4096 bytes each)
       9-12 (see below)
      13     BYTE        1     file length                length of file in sectors
      14     BYTE        1     start sector (0..15)       represent the start sector, calculated as
                                                          = start track*16+start sector
      15     BYTE        1     start track

**Case 1**: parameters of Basic program

       9     WORD        2     progs+vars                 length of program + variables area
      11     WORD        2     progs                      length of program only

**Case 2**: parameters of data arrays

       9     WORD        2     param 1                    unused
      11     WORD        2     data length                length of data array

**Case 3**: parameters of code (byte array)

       9     WORD        2     start address  
      11     WORD        2     code length

**Case 4**: parameters of print file

       9     BYTE        1     extent no.                 number of part of the print file, beginning with 0
      10     BYTE        1     unused                     always #32
      11     WORD        2     print length               0..4096


## Specification (256 bytes)

    Offset | Type   | Length | Description              | Additional Information
        0    BYTE        1     end of directory           must be 0 to indicate the end of the directory
        1    BYTE[]    224     unused                     filled with zeroes
      225    BYTE        1     first free sector (0..15)  representing the logical sector number of the next free sector on the disk.
                                                          If the disk is full, the sector number is = count of sectors.
      226    BYTE        1     first free track           logical sector=first free track*16+first free sector
      227    BYTE        1     disk type                  22: double-sided, 80 tracks
                                                          23: double-sided, 40 tracks
                                                          24: single-sided, 80 tracks
                                                          25: single-sided, 40 tracks
      228    BYTE        1     file count                 0..128; the count of non-deleted files
      229    WORD        2     free sectors               number of free sectors on the disk
      231    BYTE        1     TR-DOS ID                  always 16
      232    WORD        2     unused                     filled with 0
      234    CHAR[]      9     unused                     filled with spaces (#32), or a disk protecting password, filled up with spaces (thanks to this info by Németh Zoltán Gábor)
      243    BYTE        1     unused                     filled with 0
      244    BYTE        1     deleted files              number of deleted files on the disk
      245    CHAR[]      8     disk label                 label name of the disk
      253    BYTE[]      3     unused                     filled with 0


The file bodies are stored in the data sectors, beginning at sector #16. 

Basic program and data array files have an addition at the end of the file.
The structure of the different file types is as follows:

(`flen` = File Length)

File body of Basic program files:

    Offset    | Type   | Length | Description                     | Additional Information
    +0          BYTE[]   [flen]   data of the Basic program         the pure bytes that represent the basic (and if present, variables) data
    +[flen]     CHAR[]      2     parameter 2 indicator             always #128 #170
    +[flen]+2   WORD        2     autostart line number             0..9999


File body of data array files:

    Offset    | Type   | Length | Description                     | Additional Information
    +0          BYTE[]   [flen]   data of the data array variable   the pure bytes that represent the content of the data variable
    +[flen]     CHAR[]      2     parameter 2 indicator             always #128 #170
    +[flen]+2   BYTE        1     unused   
    +[flen]+3   BYTE        1     name of variable                  bits 0..5:  1..26 meaning "a".."z"
                                                                    bits 6..7:  10b = numeric array
                                                                                11b = alphanumeric array

File body of all other files:

    Offset    | Type   | Length | Description                     | Additional Information
     +0         BYTE[]   [flen]   data of the file                  the pure bytes that represent the content of the file


Source: http://www.zx-modules.de:80/fileformats/trdformat.html
