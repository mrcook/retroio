// Package cdt implements reading Amstrad CDT (TZX) formatted files,
// as specified in the TZX specification.
// https://www.worldofspectrum.org/TZXformat.html
//
// The `.CDT` tape image file format is identical to the `.TZX` file format designed by Tomaz Kac.
// Therefore this package is a simple wrapper around the `spectrum/tzx` package.
package cdt

import (
	"fmt"

	"github.com/mrcook/retroio/spectrum/tzx"
	"github.com/mrcook/retroio/storage"
)

type CDT struct {
	*tzx.TZX
}

func New(reader *storage.Reader) *CDT {
	return &CDT{tzx.New(reader)}
}

func (d CDT) CommandDir() {
	fmt.Println("directory listing unsupported for tapes")
}
