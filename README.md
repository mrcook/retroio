# RetroIO (rio)

A command-line utility written in the Go programming language for working
with disk and cassette tape files used by emulators of home computers from
the 1980s.

One example would be `tzx` tapes from the ZX Spectrum 8-bit home computer.

_Why? Well, I wrote this program in part to better understand how tape data
is loaded in to a ZX Spectrum, but also as an experiment in working with
binary data files._


## Supported Storage Media

Read only for:

* Amstrad:      `DSK`, `CDT`
* Commodore 64: `T64`, `TAP`
* ZX Spectrum:  `TZX`, `TAP`


### BASIC Program Listing

It is possible to output the BASIC program listings from the following media:

* ZX Spectrum: `TZX` and `TAP`

Simply add the `--bas` flag when `read`ing the image.

_Please note that decoding is currently experimental and the output may not be
considered valid BASIC, and may even be garbled or missing completely._


## Installation

    $ go get -u -v github.com/mrcook/retroio/...

To install the app after manually cloning the repository you must first change to the `rio` directory:

    $ cd retroio/rio
    $ go install


## Usage

    $ rio spectrum read /path/to/tape.tzx

The media metadata and block information will be printed to the terminal.

The program will select the correct _format_ based on the file extension,
however this can be overridden with the `--format` flag.


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
