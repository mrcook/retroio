package generalized_data

// Most commonly, pilot and sync are repetitions of the same pulse, thus they are represented
// using a very simple RLE encoding structure which stores the symbol and the number of times
// it must be repeated.
// Each symbol in the data stream is represented by a string of NB bits of the block data, where
// NB = ceiling(Log2(ASD)). Thus the length of the whole data stream in bits is NB*TOTD, or in
// bytes DS=ceil(NB*TOTD/8).
type PilotRLE struct {
	Symbol          uint8  // Symbol to be represented
	RepetitionCount uint16 // Number of repetitions
}
