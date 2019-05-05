package tzx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/mrcook/tzxit/tape"
)

type Tzx struct {
	file    *tape.File
	header  Header
	archive ArchiveInfo
	blocks  []tape.Block
}

func (t *Tzx) Read() {
	if err := t.readHeader(); err != nil {
		fmt.Print(err)
		return
	}

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

		t.readBlockData(blockID)
	}
	return nil
}

func (t *Tzx) readBlockData(id byte) {
	switch id {
	case 0x10:
		ssd := &StandardSpeedData{}
		ssd.Read(t.file)
		t.blocks = append(t.blocks, ssd)
	case 0x11:
		tsd := &TurboSpeedData{}
		tsd.Read(t.file)
		t.blocks = append(t.blocks, tsd)
	case 0x12:
		pt := &PureTone{}
		pt.Read(t.file)
		t.blocks = append(t.blocks, pt)
	case 0x13:
		sop := &SequenceOfPulses{}
		sop.Read(t.file)
		t.blocks = append(t.blocks, sop)
	case 0x14:
		pd := &PureData{}
		pd.Read(t.file)
		t.blocks = append(t.blocks, pd)
	case 0x15:
		dr := &DirectRecording{}
		dr.Read(t.file)
		t.blocks = append(t.blocks, dr)
	case 0x18:
		cr := &CswRecording{}
		cr.Read(t.file)
		t.blocks = append(t.blocks, cr)
	case 0x19:
		gd := &GeneralizedData{}
		gd.Read(t.file)
		t.blocks = append(t.blocks, gd)
	case 0x20:
		pt := &PauseTapeCommand{}
		pt.Read(t.file)
		t.blocks = append(t.blocks, pt)
	case 0x21:
		gs := &GroupStart{}
		gs.Read(t.file)
		t.blocks = append(t.blocks, gs)
	case 0x22:
		ge := &GroupEnd{}
		ge.Read(t.file)
		t.blocks = append(t.blocks, ge)
	case 0x23:
		jt := &JumpTo{}
		jt.Read(t.file)
		t.blocks = append(t.blocks, jt)
	case 0x24:
		ls := &LoopStart{}
		ls.Read(t.file)
		t.blocks = append(t.blocks, ls)
	case 0x25:
		ls := &LoopEnd{}
		ls.Read(t.file)
		t.blocks = append(t.blocks, ls)
	case 0x26:
		cs := &CallSequence{}
		cs.Read(t.file)
		t.blocks = append(t.blocks, cs)
	case 0x27:
		rs := &ReturnFromSequence{}
		rs.Read(t.file)
		t.blocks = append(t.blocks, rs)
	case 0x28:
		s := &Select{}
		s.Read(t.file)
		t.blocks = append(t.blocks, s)
	case 0x2a:
		st := &StopTapeWhen48kMode{}
		st.Read(t.file)
		t.blocks = append(t.blocks, st)
	case 0x2b:
		sl := &SetSignalLevel{}
		sl.Read(t.file)
		t.blocks = append(t.blocks, sl)
	case 0x30:
		td := &TextDescription{}
		td.Read(t.file)
		t.blocks = append(t.blocks, td)
	case 0x31:
		m := &Message{}
		m.Read(t.file)
		t.blocks = append(t.blocks, m)
	case 0x32:
		ai := ArchiveInfo{}
		ai.Read(t.file)
		t.archive = ai
	case 0x33:
		ht := &HardwareType{}
		ht.Read(t.file)
		t.blocks = append(t.blocks, ht)
	case 0x35:
		ci := &CustomInfo{}
		ci.Read(t.file)
		t.blocks = append(t.blocks, ci)
	case 0x5a: // (90 dec, ASCII Letter 'Z')
		gb := &GlueBlock{}
		gb.Read(t.file)
		t.blocks = append(t.blocks, gb)
	default:
		// probably ID's 16,17,34,35,40 (HEX)
		log.Fatalf("ID 0x%02X is deprecated/not supported", id)
	}
}

func (t *Tzx) DisplayTapeMetadata() {
	fmt.Printf("Tape image TZX revision %d.%d\n", t.header.MajorVersion, t.header.MinorVersion)
	fmt.Println()
	fmt.Println(t.archive.ToString())
	for i, block := range t.blocks {
		fmt.Printf("#%d %s\n", i+1, block.ToString())
	}
}

func (t *Tzx) Open(filename string) error {
	t.file = &tape.File{}
	return t.file.Open(filename)
}

func (t *Tzx) Close() error {
	return t.file.Close()
}
