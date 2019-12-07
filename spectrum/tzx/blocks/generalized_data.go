package blocks

import (
	"fmt"
	"log"

	"retroio/storage"
)

// GeneralizedData
// ID: 19h (25d)
// This block has been specifically developed to represent an extremely wide range of data encoding techniques.
// The basic idea is that each loading component (pilot tone, sync pulses, data) is associated to a
// specific sequence of pulses, where each sequence (wave) can contain a different number of pulses
// from the others. In this way we can have a situation where bit 0 is represented with 4 pulses and
// bit 1 with 8 pulses.
type GeneralizedData struct {
	Length       uint32     // DWORD Block length (without these four bytes)
	Pause        uint16     // WORD  Pause after this block (ms)
	TOTP         uint32     // DWORD Total number of symbols in pilot/sync block (can be 0)
	NPP          uint8      // BYTE  Maximum number of pulses per pilot/sync symbol
	ASP          uint8      // BYTE  Number of pilot/sync symbols in the alphabet table (0=256)
	TOTD         uint32     // DWORD Total number of symbols in data stream (can be 0)
	NPD          uint8      // BYTE  Maximum number of pulses per data symbol
	ASD          uint8      // BYTE  Number of data symbols in the alphabet table (0=256)
	PilotSymbols []Symbol   // 0x12  SYMDEF[ASP] Pilot and sync symbols definition table: this field is present only if TOTP>0
	PilotStreams []PilotRLE // 0x12+ (2*NPP+1)*ASP - PRLE[TOTP]  Pilot and sync data stream: this field is present only if TOTP>0
	DataSymbols  []Symbol   // 0x12+ (TOTP>0)*((2*NPP+1)*ASP)+TOTP*3  - SYMDEF[ASD] Data symbols definition table: this field is present only if TOTD>0
	DataStreams  []uint8    // 0x12+ (TOTP>0)*((2*NPP+1)*ASP)+ TOTP*3+(2*NPD+1)*ASD - BYTE[DS]  Data stream: this field is present only if TOTD>0
}

// The alphabet is stored using a table where each symbol is a row of pulses. The number of columns
// (i.e. pulses) of the table is the length of the longest sequence amongst all (MAXP=NPP or NPD,
// for pilot/sync or data blocks respectively); shorter waves are terminated by a zero-length pulse
// in the sequence.
// Any number of data symbols is allowed, so we can have more than two distinct waves; for
// example, imagine a loader which writes two bits at a time by encoding them with four distinct
// pulse lengths: this loader would have an alphabet of four symbols, each associated to a specific
// sequence of pulses (wave).
type Symbol struct {
	Flags uint8 // BYTE Symbol flags:
	//                    b0-b1: starting symbol polarity
	//                    00:    opposite to the current level (make an edge, as usual) - default
	//                    01:    same as the current level (no edge - prolongs the previous pulse)
	//                    10:    force low level
	//                    11:    force high level
	PulseLengths []uint16 // WORD[MAXP] Array of pulse lengths.
}

// Most commonly, pilot and sync are repetitions of the same pulse, thus they are represented
// using a very simple RLE encoding structure which stores the symbol and the number of times
// it must be repeated.
// Each symbol in the data stream is represented by a string of NB bits of the block data, where
// NB = ceiling(Log2(ASD)). Thus the length of the whole data stream in bits is NB*TOTD, or in
// bytes DS=ceil(NB*TOTD/8).
type PilotRLE struct {
	Symbol          uint8  // BYTE  Symbol to be represented
	RepetitionCount uint16 // WORD  Number of repetitions
}

// Read the tape and extract the data.
// It is expected that the tape pointer is at the correct position for reading.
func (g *GeneralizedData) Read(reader *storage.Reader) {
	log.Fatal("GeneralizedData is not currently supported - unable to continue.")
}

// Id of the block as given in the TZX specification, written as a hexadecimal number.
func (g GeneralizedData) Id() uint8 {
	return 0x19
}

// Name of the block as given in the TZX specification.
func (g GeneralizedData) Name() string {
	return "Generalized Data"
}

// ToString returns a human readable string of the block data
func (g GeneralizedData) ToString() string {
	return fmt.Sprintf("%s", g.Name())
}
