# RetroIO (rio)

A command-line utility written in the Go programming language for working
with disk and cassette tape files used by emulators of home computers from
the 1980s.

One example would be `tzx` tapes from the ZX Spectrum 8-bit home computer.


## Supported Media and Commands


### Directory Command

* Amstrad:      `DSK`
* Commodore 64: `D64`, `D71`, `D81`

The `dir` command reads a disk and prints the directory listing to the terminal.
Any hidden/scratch files will also be displayed.

```sh
$ rio c64 dir super-mario-bros64.d64

LOAD"$",8
SEARCHING FOR $
LOADING
READY.
LIST

0 "SUPER MARIO BROS" .  4 
192  "SUPER M. BROS.64" PRG
59   "SMB.64 DOCS"      PRG
0    "                " DEL
0    "     ¦¦¦¦¦      " DEL
0    "    ¦¦¦¦¦¦¦¦¦   " DEL
0    "    ---vv-v     " DEL
0    "   -v-vvv-vvv   " DEL
0    "   -v--vvv-vvv  " DEL
0    "   --vvvv----   " DEL
0    "     vvvvvvv    " DEL
0    "    --¦---      " DEL
0    "   ---¦--¦---   " DEL
0    "  ----¦¦¦¦----  " DEL
0    "  vv-¦v¦¦v¦-vv  " DEL
0    "  vvv¦¦¦¦¦¦vvv  " DEL
0    "  vv¦¦¦¦¦¦¦¦vv  " DEL
0    "    ¦¦¦  ¦¦¦    " DEL
0    "   ---    ---   " DEL
0    "  ----    ----  " DEL
0    "                " DEL
428 BLOCKS FREE.
```


### Geometry Command

* Amstrad:      `DSK`, `CDT`
* Commodore 64: `D64`, `D71`, `D81`, `T64`, `TAP`
* ZX Spectrum:  `TZX`, `TAP`

The `geometry` command will read and display core metadata about the layout
of the media. This can be disk track and sector details, or the header and
block information from a cassette tape.

```sh
$ rio spectrum geometry skool-daze.tzx

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


### Read Command

* ZX Spectrum: `TZX` and `TAP`

The `read` command will read data contained on the media.

At present only printing of `BASIC` programs is supported. Simply add the `--bas`
flag when `read`ing the media image.

_Please note that decoding is currently experimental and the output may not be
considered valid BASIC, and may even be garbled or missing completely._

```sh
$ rio spectrum read --bas manic-miner.tzx

BASIC PROGRAMS:

BLK#02: ManicMiner
  10  CLEAR 30000
  20  PAPER 0: BORDER 0: INK 0: CLS : LOAD ""CODE : LOAD ""CODE 
  30  RANDOMIZE USR 33792
```

## Installation

    $ go get -u -v github.com/mrcook/retroio/...

To install the app after manually cloning the repository you must first change to the `rio` directory:

    $ cd retroio/rio
    $ go install


## Usage

    $ rio help

To display a list of commands

Note: the correct media type will be set for the requested system based
on the file extension, however this can be overridden with the `--media` flag.


## LICENSE

Copyright (c) 2018-2020 Michael R. Cook. All rights reserved.

This work is licensed under the terms of the MIT license.
For a copy, see <https://opensource.org/licenses/MIT>.
