package tape

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type Tzx struct {
	file   *File
	header Header
	blocks []Block
}

func (t *Tzx) Process() {
	if err := t.readHeader(); err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("TZX revision %d.%d\n", t.header.MajorVersion, t.header.MinorVersion)

	if err := t.readBlocks(); err != nil {
		fmt.Print(err)
		return
	}
}

func (t *Tzx) readHeader() error {
	t.header = Header{}
	data := t.file.ReadBytes(10)

	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.LittleEndian, &t.header)
	if err != nil {
		return fmt.Errorf("binary.Read failed: %v", err)
	}

	if string(t.header.Signature[:]) != "ZXTape!" {
		return fmt.Errorf("TZX file is not in correct format")
	}

	return nil
}

// readBlocks processes all the TZX data blocks
func (t *Tzx) readBlocks() error {
	for {
		blockID, err := t.file.ReadByte()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		t.processBlockData(blockID)
		break // because processBlockData() is not implemented yet
	}
	return nil
}

func (t *Tzx) processBlockData(id byte) {
	var block Block

	switch id {
	case 16:
		ssd := &StandardSpeedData{}
		ssd.Process(t.file)
		t.blocks = append(t.blocks, ssd)
	case 17:
		tsd := &TurboSpeedData{}
		tsd.Process(t.file)
		t.blocks = append(t.blocks, tsd)
	case 18:
		pt := &PureTone{}
		pt.Process(t.file)
		t.blocks = append(t.blocks, pt)
	case 19:
		sop := &SequenceOfPulses{}
		sop.Process(t.file)
		t.blocks = append(t.blocks, sop)
	case 20:
		pd := &PureData{}
		pd.Process(t.file)
		t.blocks = append(t.blocks, pd)
	case 21:
		dr := &DirectRecording{}
		dr.Process(t.file)
		t.blocks = append(t.blocks, dr)
	case 24: // CswRecording()
	case 25: // GeneralizedData()
	case 32: // PauseTheTapeCommand()
	case 33: // GroupStart()
	case 34: // GroupEnd()
		return // block has no body
	case 35: // JumpTo()
	case 36: // LoopStart()
	case 37: // LoopEnd()
		return // block has no body
	case 38: // CallSequence()
	case 39: // ReturnFromSequence()
		return // block has no body
	case 40: // Select()
	case 42: // StopTheTapeIfIn48kMode()
	case 43: // SetSignalLevel()
	case 48: // TextDescription()
	case 49: // Message()
	case 50: // ArchiveInfo()
	case 51: // HardwareType()
	case 53: // CustomInfo()
	case 90: // GlueBlock()
	default:
		// probably ID's 16,17,34,35,40 (HEX) / 22,23,52,64 (DECIMAL)
		log.Fatalf("ID %d is deprecated/not supported", id)
	}

	t.blocks = append(t.blocks, block)
}

func (t *Tzx) Open(filename string) error {
	t.file = &File{}
	return t.file.Open(filename)
}

func (t *Tzx) Close() error {
	return t.file.Close()
}
