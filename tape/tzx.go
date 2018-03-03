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
	Head   Header
	Blocks []Block
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

func (t *Tzx) Run() {
	t.Head = Header{}
	data := t.readNextBytes(10)

	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.LittleEndian, &t.Head)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	if string(t.Head.Signature[:]) != "ZXTape!" {
		log.Fatal("given TZX file is not in correct format")
	}

	fmt.Printf("TZX v%d.%d\n", t.Head.MajorVersion, t.Head.MinorVersion)
}

func (t Tzx) readNextBytes(number int) []byte {
	b := make([]byte, number)
	_, err := t.file.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
