# Atari CAS / A8CAS file format

Source: http://a8cas.sourceforge.net/format-cas.html

CAS is a binary file format originally introduced by Ernest R. Schreurs' as the
output format of his WAV2CAS utility. It is intended to store Atari tapes
efficiently.

Originally, CAS format could not store non-standard signals that are common on
tapes with copy-protection measures. In A8CAS, the file format has been
extended to support those tapes, and also to support tapes in turbo formats.
The extension is backwards-compatible - that is, all "old" CAS files are
readable by `liba8cas`.


## Structure

The file in CAS format is composed of chunks, each having a common structure:
an 8-byte header and some data.

Each chunk consists of:

    Offset | Size (bytes)  | Name         | Contains
       0      4              chunk_type     4-letter string.
       4      2              chunk_length   Chunk length (not including header)
       6      2              aux            Contents depend on chunk type
       8     chunk_length    data           Contents depend on chunk type

So, each chunk has size `chunk_length+8` bytes.

All numeric values stored in a CAS file are **little-endian**.


## Chunk types

There are several possible chunk types, each identified by a 4-letter string.
Those are:


### FUJI - tape description

    Offset | Size         | Name         | Contains
       0      4             chunk_type     FUJI
       4      2             chunk_length   chunk's length
       6      2                            ignored
       8     chunk_length   description    tape's description (UTF-8)

A CAS file must begin with a `FUJI` chunk. This chunk holds the tape's description
(as an UTF-8 string) in data. The `chunk_length` field holds the description's
length. The description may be empty - then the `chunk_length` is `0`.

A CAS file may contain more FUJI chunks - they can be used to label different
parts of a tape, such as separate files recorded one after another. However the
current version of `liba8cas` ignores all but the first FUJI chunk.


### baud - baudrate for subsequent SIO records

    Offset | Size | Name         | Contains
       0      4     chunk_type     baud
       4      2     chunk_length   0x00 00 (0 bytes)
       6      2     baudrate       baudrate of subsequent data chunks

A baud chunk holds the baud rate of the following data chunks (until a next baud
chunk). The chunk's length is always 0 (+8 bytes of the header), and the
baud rate is stored in the `baudrate` field. If no baud chunk is encountered
before a data chunk, its baud rate is set to the default value of 600.


### data - standard SIO record

    Offset | Size         | Name         | Contains
       0      4             chunk_type     data
       4      2             chunk_length   chunk's length
       6      2             irg_length     length of IRG before this record, in ms
       8     chunk_length   data           block's data

A data chunk contains a standard tape record as read or written by Atari's SIO.
Those records normally have a length of 132 bytes, start with two `0x55` bytes,
and end with a `checksum` byte - however those are not necessary. The record's
baud rate is the `baudrate` stored in the previous baud chunk (if no baud chunk
has been encountered yet, a standard baud rate of 600 is assumed).

The `irg_length` field holds length of an Inter-Record Gap (IRG) before this
chunk, in milliseconds. The `chunk_length` field contains number of bytes in
the record. `data` contains all the record's bytes.


### fsk  - non-standard SIO signals

    Offset | Size         | Name         | Contains
       0      4             chunk_type     fsk  (ends with space - the length is 4 characters)
       4      2             chunk_length   chunk's length
       6      2             irg_length     length of IRG before this record, in ms
       8     chunk_length   data           lengths (in 1/10's ms) of alternating 0/1 signals, as unsigned 16-bit values

An fsk  chunk holds a raw sequence of lengths of FSK signals (0s and 1s). It is
needed to store special signals written on a tape, or any data that could not
be recognised as a standard tape record.

The `irg_length` field holds length of an Inter-Record Gap (IRG) before this
chunk, in milliseconds. The `data` field consists of `chunk_length/2` 16-bit
unsigned numbers (LSB first) - each number represents length of a single
signal (0 or 1) in 1/10's of milliseconds.

Example: the sequence
    
    0x66 73 6b 20 0a 00 11 01 00 01 10 01 80 00 20 00 80 02

is an `fsk` block that means:

    Name         | Raw bytes   | value
    chunk_type     66 73 6b 20   chunk's identifier
    chunk_length   0a 00         Length is 0x000a = 10 bytes
    aux            11 01         IRG is 0x0111 = 273 milliseconds
                   00 01         0x0100 = 256: 25.6ms of logical 0
                   10 01         0x0110 = 272: 27.2ms of logical 1
    data           80 00         0x0080 = 128: 12.8ms of logical 0
                   20 00         0x0020 = 32: 3.2ms of logical 1
                   80 02         0x0280 = 640: 64ms of logical 0


### pwms - settings for subsequent turbo records

    Offset | Size | Name         | Contains
      0       4     chunk_type     pwms
      4       2     chunk_length   0x02 00 (2 bytes)
      6       1     ...            see description below
      7       1     ignored        reserved for future use
      8       2     samplerate     a base value in Hz, from which subsequent pulse lengths will be derived

Offset 6, size: 1-byte

    Bits | Name       | Contains
    0-1    pulse_type   %01 - a pulse consists of a falling edge and then a rising edge in that order. Both edges have equal length.
                        %10 - a pulse consists of a rising edge and then a falling edge in that order. Both edges have equal length.
     2     bit_order    0 - least significant bit first
                        1 - most significant bit first
    3-7    ignored      reserved for future use


A pwms chunk has similar meaning for turbo transmission as a baud chunk for
normal SIO transmission. It defines parameters for decoding all subsequent
pwmc, pwmd and pwml chunks (until a next pwms chunk). The chunk's length is
always 2 (+8 bytes of the header).

In turbo transmission, bits are encoded using PWM - bits `0` and `1` are
represented by pulses of different lengths. A single pulse consists of two
edges: one rising, one falling. The `pulse_type` field defines the exact
structure of a pulse. (In most turbo systems, a rising edge in a sound signal
causes a logical 1 to appear on SIO's DATA IN pin; a falling edge causes a
logical 0 to appear.)

The `bit_order` field defines bit order of data stored in subsequent pwmd chunks.

The `samplerate` field is used in pwmc, pwmd and pwml chunks to determine lengths
of pulses. If a pulse has given length of n, it means that its real length is
`n/samlperate` seconds.


### pwmc - sequence of turbo signals

    Offset | Size         | Name           | Contains
      0       4             chunk_type       pwmc
      4       2             chunk_length     chunk's length
      6       2             silence_length   length of silence before this block, in ms
      8      chunk_length   data             sequence of 3-byte elements, each defining length of a single signal:
                                               Size | Contains
                                                1     Length of a single pulse
                                                2     Length of the signal, in pulses

A pwmc chunk holds a sequence of signals. A signal is simply a sinewave of
specified frequency, lasting a specified amount of pulses. Such signals appear
at the beginning of each turbo data block and are used for synchronisation.

The `data` field contains signals, each one defined by 3 bytes. Therefore a pwmc
chunk holds `chunk_length/3` signals.

Each signal is defined in terms of 2 values: length of a single pulse and
number of pulses. For example, a signal defined by bytes `0x03 20 01` is a signal
that consists of `0x0120` (276) pulses each of length `0x03` (3) (total signal's
length = 276x3. See `samplerate` field in the last encountered pwms chunk for
length's base rate).


### pwmd - turbo record with data

    Offset | Size         | Name           | Contains
      0       4             chunk            type  pwmd
      4       2             chunk_length     chunk's length
      6       1             pulse_0_length   length of pulses representing logical 0
      7       1             pulse_1_length   length of pulses representing logical 1
      8      chunk_length   data             block's data

A pwmd chunk has similar purpose as a data chunk, but for turbo transmission - it
holds a block of bytes encoded in PWM. The fields `pulse_0_length` and
`pulse_1_length` contain lengths of pulses that are used do encode bits in this
chunk. The `data` field contains the block's bytes. Bit order of the block's bytes
is defined by the `bit_order` field in the last encountered pwms chunk.


### pwml

    Offset | Size         | Name           | Contains
       0      4             chunk            type  pwml
       4      2             chunk_length     chunk's length
       6      2             silence_length   length of silence before this block, in ms
       8     chunk_length   data             lengths of PWM states, as unsigned 16-bit values.

A pwml chunk has similar purpose as a fsk  chunk, but for turbo transmission - it
holds a sequence of raw PWM states. Each state can be high = rising edge, or
low = falling edge. pwml chunks are used as a last-resort method of storing PWM
data that cannot be encoded into bytes (in other words, into a pwmd chunk).

A pwml chunk contains `chunk_length/2` states - each state is encoded as a 2-byte
number in the `data` field. Like always, a signals length should be divided by
samplerate from the previous pwms chunk to get the length in seconds.

Interpretation of the values in the `data` field depends on the previous pwms
block's `pulse_type` value:

* If `pulse_type` equals `%10`, then the first value in data is length of a rising
  edge, the second value is length of a falling edge, the third value - again
  a rising edge, and so on.
* If `pulse_type` equals `%01`, then the first value in data is length of a falling
  edge, the second value is length of a rising edge, the third value - again a
  falling edge, and so on.
