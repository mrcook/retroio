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
	case 24:
		cr := &CswRecording{}
		cr.Process(t.file)
		t.blocks = append(t.blocks, cr)
	case 25:
		gd := &GeneralizedData{}
		gd.Process(t.file)
		t.blocks = append(t.blocks, gd)
	case 32:
		pt := &PauseTapeCommand{}
		pt.Process(t.file)
		t.blocks = append(t.blocks, pt)
	case 33:
		gs := &GroupStart{}
		gs.Process(t.file)
		t.blocks = append(t.blocks, gs)
	case 34:
		ge := &GroupEnd{}
		ge.Process(t.file)
		t.blocks = append(t.blocks, ge)
	case 35:
		jt := &JumpTo{}
		jt.Process(t.file)
		t.blocks = append(t.blocks, jt)
	case 36:
		ls := &LoopStart{}
		ls.Process(t.file)
		t.blocks = append(t.blocks, ls)
	case 37:
		ls := &LoopEnd{}
		ls.Process(t.file)
		t.blocks = append(t.blocks, ls)
	case 38:
		cs := &CallSequence{}
		cs.Process(t.file)
		t.blocks = append(t.blocks, cs)
	case 39:
		rs := &ReturnFromSequence{}
		rs.Process(t.file)
		t.blocks = append(t.blocks, rs)
	case 40:
		s := &Select{}
		s.Process(t.file)
		t.blocks = append(t.blocks, s)
	case 42:
		st := &StopTapeWhen48kMode{}
		st.Process(t.file)
		t.blocks = append(t.blocks, st)
	case 43:
		sl := &SetSignalLevel{}
		sl.Process(t.file)
		t.blocks = append(t.blocks, sl)
	case 48:
		td := &TextDescription{}
		td.Process(t.file)
		t.blocks = append(t.blocks, td)
	case 49:
		m := &Message{}
		m.Process(t.file)
		t.blocks = append(t.blocks, m)
	case 50:
		ai := &ArchiveInfo{}
		ai.Process(t.file)
		t.blocks = append(t.blocks, ai)
	case 51:
		ht := &HardwareType{}
		ht.Process(t.file)
		t.blocks = append(t.blocks, ht)
	case 53:
		ci := &CustomInfo{}
		ci.Process(t.file)
		t.blocks = append(t.blocks, ci)
	case 90:
		gb := &GlueBlock{}
		gb.Process(t.file)
		t.blocks = append(t.blocks, gb)
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
