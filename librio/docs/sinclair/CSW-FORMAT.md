# Compressed Square Wave (CSW) Specification

Created by Ramsoft ZX Spectrum demogroup.

Format revision: v2.00 (August 1st 2003).

* Introduction
* The CSW utility
* CSW file format
* Old CSW 1.01 file format
* Contact information
* Revision history


## Introduction

This document describes the CSW file format and the CSW.EXE utility. CSW is
strongly based upon the MakeTZX engine and it shares with it various aspects
of its behaviour. In the manual of MakeTZX you will find lots of explanations,
tips and FAQ that are not reported here, so we recommend you to read it too.


### The CSW file format

CSW files are a way of storing sample data in a compact form, typically taking
1/10th of an ordinary VOC. It is used internally by MakeTZX, but it is also very
useful to keep down the disk space taken by your VOC/WAV files. The CSW utility
can handle CSW conversion in both ways (see below). Of course, MakeTZX itself
accepts CSW files for input. When converting to the CSW format, the sample file
is processed through MakeTZX's internal digital filter which reduces noise and
signal distortions very efficiently. Make a backup copy of the original file if
you will need the original samples later, but remember that in most cases the
CSW will be a lot better than the original file. Note that CSWs are intended
for use with square waves only (such as computer tapes)! The compression ratio
depends on many factors; in general, the higher the sample rate, the higher the
ratio. A clean and regular signal helps too. The ratio for a 44 KHz file will
usually be twice the value for a 22 KHz one. The typical gain for a 44 KHz
turbo tape is about 93%, which means a 12:1 compression factor! Normal speed
tapes should compress even better. Finally, CSW files are highly compressible
with the standard PC archivers such as RAR and ZIP. The packed CSW files are
usually smaller than the zipped original VOCs. You will be able to RAR a 40 MB
sample file down to a few hundreds KB.


## The CSW utility

This small program is intended to provide a basic support for CSW files. It can
compress VOC, WAV, IFF and OUT files to CSW and decompress CSW files back to VOC
format (switch -d). Enter CSW -? or simply CSW for help. At the moment, CSW.EXE
accepts only uncompressed mono 8-bit sample files. Extensions in filenames can
be omitted; in this case, the default extensions will be appended in turn to
match an existing file. The search order is VOC, WAV, IFF and OUT for last. If
the output filename is left out, the input file name with extension .CSW (or
.VOC if decompressing) will be used. If the input filename ends with the
extension .CSW, then the switch -d (decompression) is implicitly assumed.

CSW can also work in **DirectMode** (switch -r), in which case the input is
taken from your sound card and conversion is performed on-the-fly in true
realtime. You can stop the conversion by pressing any key at any time. To pause
the recording press 'P', followed by any key to resume. During the pause, the
vu-meter is shown again. Note that, due to MakeTZX's engine requirements, the
samples are written to disk anyway, so the maximum recording time is limited by
the available disk space. If you want, you can keep this samples at the end of
conversion and save them in a WAV file (switch -k), just in case something goes
wrong and you don't want to repeat the sample. In this way, CSW may also act
like a sampler! You can set the sampling frequency with switch -s (e.g. -s44100).
You can also do programmed recordings using switch -t and specifying the
recording time (in seconds, e.g. -t60.0 for one minute); in this case, CSW will
automatically stop when the time has elapsed (or when the disk is full), so you
can start it and go away to do better things :)

The DirectMode SoundBlaster driver has been written for 100% compatible
sound cards. If you are experiencing problems, try option "-c" which will
attempt to access the hardware in a different way. The driver also performs a
preliminary stability check; if this fails, CSW exits after two seconds with an
error message. All this stuff is extensively covered in MakeTZX's manual,
DirectMode section; please read it carefully. Note: In order to run the CSW
utility under plain MS-DOS you need a DPMI host (such as CWSDPMI.EXE).

Note: DirectMode, OUT files, digital filter and the other features are
extensively described into MakeTZX's manual. Please read it.

Note: Although it is possible to specify fractions of seconds, the effective
recording time is subject to DMA buffer size quantums (a few 1/10ths of sec).

Note: Like MakeTZX, CSW supports long filenames under Windows 9x


## CSW-2 file format

Here is the CSW implementation chart for anyone who wants to use it in some
utility or emulator (if so, please let us know). The file format is very simple
and the compression scheme used is somewhat based on the RLE algorithm.

    Legenda

    WORD       2 bytes
    DWORD      4 bytes
    BYTE[N]    N bytes
    ASCII[N]   N ASCII characters
    ASCIIZ[N]  ASCII string with zero-padding to N bytes total

* All multi-byte values are stored in Intel byte order (little-endian).
* All reserved or undefined bits must be set to zero.
* All the headers fields must be filled in; blank values are not allowed.


    CSW-2 Header

    CSW global file header - status: required

    Offset   | Value  | Type       | Description
    ---------|--------|------------|-----------------------------------------
    0x00     | (note) | ASCII[22]  | "Compressed Square Wave" signature
    0x16     | 0x1A   | BYTE       | Terminator code
    0x17     | 0x02   | BYTE       | CSW major revision number
    0x18     | 0x00   | BYTE       | CSW minor revision number
    0x19     |  -     | DWORD      | Sample rate
    0x1D     |  -     | DWORD      | Total number of pulses (after decompression)
    0x21     |  -     | BYTE       | Compression type (see notes below)
    0x22     |  -     | BYTE       | Flags. b0: initial polarity; if set, the signal starts at logical high
    0x23     | HDR    | BYTE       | Header extension length in bytes (0x00). For future expansions only, see note below.
    0x24     |  -     | ASCIIZ[16] | Encoding application description. Information about the tool which created the file (e.g. name and version)
    0x34     |  -     | BYTE[HDR]  | Header extension data (if present)
    0x34+HDR |  -     |    -       | CSW data

### Note about Header Extensions

CSW-2 allows to extend the header size by a certain amount of bytes (the
current default value is 0). However, this is designed for future revisions of
this format and it is not meant to store application-specific data.

### Compression types

* `0x01`: RLE (Run Length Encoding)
  - The data is stored as a sequence of pulse lengths (1 byte per pulse).
    Consider the following scenario (each dot is a sample):
  - The 5 pulses shown will be represented with the following bytes: `03 05 01 04 07`
  - Pulse lengths greater than `0xFF` (255) are represented as byte `0x00`
    followed by the duration represented on 4 bytes, e.g. `0xCDE9` is stored
    as `00 E9 CD 00 00`.
* `0x02`: Z-RLE (CSW v2.xx only)
  - Pulses are encoded exactly as in method 1, but the generated byte-stream is
    further compressed with the standard `deflate()` algorithm as defined by
    the ZLIB library (RFC 1151 and 1152). In fact the compression is equivalent
    to `gzip -9` (without the magic signature); the source code of the
    compression routines we used is the same as in our RZX SDK.

In format revision 1.01 we have introduced a bit to represent the initial
signal polarity, which is not important in the Spectrum world but it is for
other platforms such as C64. All the Spectrum TZX converters can safely ignore
this bit (like MakeTZX does), so any tool supporting CSW 1.00 will also work
fine with CSW 1.01 without modifications.

Note that no info about the pulse amplitude is represented because it is not
necessary, since we are dealing with discrete 2-values amplitude scales.


## Old CSW v1.01 file format

This is the format specification for the old CSW v1.01. It is reported here
because a lot of existing tools support the original version of the file format.

    CSW-1 Header

    CSW global file header - status: required

    Offset   | Value  | Type       | Description
    ---------|--------|------------|-----------------------------------------
    0x00     | (note) | ASCII[22]  | "Compressed Square Wave" signature
    0x16     | 0x1A   | BYTE       | Terminator code
    0x17     | 0x02   | BYTE       | CSW major revision number
    0x18     | 0x00   | BYTE       | CSW minor revision number
    0x19     |   -    | WORD       | Sample rate
    0x1B     | 0x01   | BYTE       | Compression type (see notes below). 0x01: RLE.
    0x1C     |   -    | BYTE       | Flags. b0: initial polarity; if set, the signal starts at logical high
    0x1D     | 0x00   | BYTE[3]    | Reserved.
    0x20     |   -    |   -        | CSW data.

For information about the RLE compression method (0x01) and the meaning of the
polarity flag, please refer to the notes for version 2.xx.


## Contact information

The latest version of this document can be found at: http://www.ramsoft.bbk.org/csw.html

E-mails concerning the CSW specifications should be directed to: ramsoft@bbk.org


## Revision history

* Revision 2.00 (August 1st 2003)
 - Introduced CSW revision 2.00.
 - Cleared up the document a bit.
* Revision 1.01 (July 13th 1999)
 - Introduced the polarity bit (b0 in Flags)
 