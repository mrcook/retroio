# RetroIO (rio)

A command-line utility written in the Go programming language for working
with disk and cassette tape files used by emulators of home computers from
the 1980s.

One example would be `tzx` tapes from the ZX Spectrum 8-bit home computer.


## Supported Media and Commands


### Geometry Command

* Amstrad:      `DSK`, `CDT`
* Commodore 64: `T64`, `TAP`
* ZX Spectrum:  `TZX`, `TAP`

The `geometry` command will read and display core metadata about the layout
of the media. This can be disk track and sector details, or the header and
block information from a cassette tape.


### Directory Command

* Amstrad:      `DSK`

The `dir` command reads a disk and prints the directory listing to the terminal.
Any hidden files will also be displayed.


### Read Command

* ZX Spectrum: `TZX` and `TAP`

The `read` command will read data contained on the media.

At present only printing of `BASIC` programs is supported. Simply add the `--bas`
flag when `read`ing the media image.

_Please note that decoding is currently experimental and the output may not be
considered valid BASIC, and may even be garbled or missing completely._


## Installation

    $ go get -u -v github.com/mrcook/retroio/...

To install the app after manually cloning the repository you must first change to the `rio` directory:

    $ cd retroio/rio
    $ go install


## Usage

    $ rio spectrum geometry /path/to/tape.tzx

Reads the media geometry and prints the details to the terminal.

The program will select the correct media type for the requested system, based
on the file extension, however this can be overridden with the `--media` flag.


### Example output

```
TZX processing complete!

ARCHIVE INFORMATION:
  Title     : Skool Daze
  Publisher : Microsphere
  Authors   : David S. Reidy, Keith Warrington
  Year      : 1984
  Loader    : Microsphere
  Comment   : Timing corrected by Mikie.

DATA BLOCKS:
#1 Standard Speed Data : 19 bytes, pause for 970 ms.
   - Header       : BASIC Program
   - Filename     : skooldaze
   - AutoStartLine: 0
#2 Standard Speed Data : 333 bytes, pause for 5981 ms.
   - Standard Data: 331 bytes
#3 Turbo Speed Data    : 82109 bytes, pause for 0 ms.

TZX revision: 1.10
```


## LICENSE

Copyright (c) 2018-2020 Michael R. Cook. All rights reserved.

This work is licensed under the terms of the MIT license.
For a copy, see <https://opensource.org/licenses/MIT>.
