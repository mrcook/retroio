package tape

import (
	"log"
	"os"
)

type File struct {
	file *os.File
}

func (t *File) Open(filename string) error {
	var err error
	t.file, err = os.Open(filename)
	if err != nil {
		return err
	}
	return nil
}

func (t File) Close() error {
	return t.file.Close()
}

// readNextByte reads one byte from the file, usually this is the block ID byte.
func (t File) ReadNextByte() (byte, error) {
	var b [1]byte
	_, err := t.file.Read(b[:])
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

// readBlockData is used to read the tape block data.
// Note: we should never reach EOF with this read - if it happens, then there's
// something seriously wrong with the tape image - so exit with a fatal.
func (t File) ReadBlockData(numberBytes int) []byte {
	b := make([]byte, numberBytes)
	_, err := t.file.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}
