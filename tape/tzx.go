package tape

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type Tzx struct {
	file   *os.File
	header Header
	blocks []Block
}

func (t *Tzx) Process() {
	if err := t.readHeader(); err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("TZX revision %d.%d\n", t.header.MajorVersion, t.header.MinorVersion)
}

func (t *Tzx) readHeader() error {
	t.header = Header{}
	data := t.readNextBytes(10)

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

func (t Tzx) readNextBytes(number int) []byte {
	b := make([]byte, number)
	_, err := t.file.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func (t *Tzx) Open(filename string) error {
	var err error
	t.file, err = os.Open(filename)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tzx) Close() error {
	return t.file.Close()
}
