# TZX Utility

A utility program for working with the ZX Spectrum tape files (TZX), written
in the Go language.

I wrote this program, in part to better understand how tape data is loaded in
to a spectrum, but also as an experiment in working with binary data files.

## Installation

    $ go get -u -v github.com/mrcook/tzxit/...

## Usage

    $ tzxit read /path/to/tape.tzx

The TZX metadata will be printed to the terminal.

### Example output

```
TZX Revision: 1.10

  Title     : Skool Daze
  Publisher : Microsphere
  Authors   : David S. Reidy, Keith Warrington
  Year      : 1984
  Loader    : Microsphere
  Comment   : Timing corrected by Mikie.

> Standard Speed Data : 19 bytes, pause for 970 ms.
  - Header       : BASIC Program
  - Filename     : skooldaze 
  - AutoStartLine: 0
> Standard Speed Data : 333 bytes, pause for 5981 ms.
> Turbo Speed Data    : 82109 bytes, pause for 0 ms.
```

## TZX Specification

Sources:

- https://www.worldofspectrum.org/TZXformat.html
- http://www.zx-modules.de/fileformats/tapformat.html

`TZX` is a file format designed to preserve cassette tapes compatible with the
ZX Spectrum computers, although some specialized versions of the format have
been defined for other machines such as the Amstrad CPC and C64.

The format was originally created by Tomaz Kac, who was maintainer until
`revision 1.13`, before passing it to Martijn v.d. Heide. For a brief period
the Ramsoft company became the maintainers, and created revision `v1.20`.

The default file extension is `.tzx`.

The tape files processable with this program are based on the TZX specification,
revision: 1.20 (2006-12-19), therefore the following hex block ID's are not
supported: `16`, `17`, `34`, `35`, and `40`.

NOTE: `GeneralizedData` blocks are also not currently supported.


## LICENSE

Copyright (c) 2018-2019 Michael R. Cook. All rights reserved.

This work is licensed under the terms of the MIT license.
For a copy, see <https://opensource.org/licenses/MIT>.
