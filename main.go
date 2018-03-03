package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/mrcook/spectrumator/tzx"
)

func main() {
	filename := "tape-file.tzx"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	defer file.Close()

	fmt.Printf("%s opened successfully!\n", filename)

	header := tzx.Header{}
	data := readNextBytes(file, 10)

	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.LittleEndian, &header)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	if string(header.Signature[:]) != "ZXTape!" {
		log.Fatal("given TZX file is not in correct format")
	}

	fmt.Printf("TZX v%d.%d\n", header.MajorVersion, header.MinorVersion)
}

func readNextBytes(file *os.File, number int) []byte {
	b := make([]byte, number)
	_, err := file.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
