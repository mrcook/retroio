// https://faqwiki.zxnet.co.uk/wiki/TR-DOS_filesystem
// http://www.8bit-wiki.de/index.php?id=3&filt=Sinclair/ZX_Spectrum/interfaces/betadisk/_manual/&cid=16062&mode=dl&tx=330ea210000
package trd

import (
	"fmt"
	"github.com/mrcook/retroio/storage"
)

type TRD struct {
	reader *storage.Reader

	Files []FileInformation
	Info  DiskInformation
}

func New(reader *storage.Reader) *TRD {
	return &TRD{reader: reader}
}

// Read processes each TAP/BLK block in the tape file.
func (t *TRD) Read() error {
	var descriptors [128]FileInformation

	for i := 0; i < len(descriptors); i++ {
		err := descriptors[i].Read(t.reader)

		if err != nil {
			return err
		}
	}

	info := DiskInformation{}
	err := info.Read(t.reader)

	if err != nil {
		return err
	}

	t.Info = info

	t.Files = descriptors[:info.NumFiles]

	return nil
}

// DisplayGeometry outputs the metadata of the disk and its files to the terminal.
func (t TRD) DisplayGeometry() {
	fmt.Println("DISK INFORMATION:")
	fmt.Println(t.Info)

	fmt.Println("FILES:")
	for i, file := range t.Files {
		if file.IsDeleted() {
			continue
		}
		fmt.Printf("#%03d %s\n", i+1, file.String())
	}
}

// DisplayBASIC outputs all BASIC programs on the disk
func (t TRD) DisplayBASIC() {
	fmt.Println("Not implemented for Spectrum disk images")
}
