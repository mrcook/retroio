// Package cdt implements reading Amstrad CDT (TZX) formatted files,
// as specified in the TZX specification.
// https://www.worldofspectrum.org/TZXformat.html
//
// The `.CDT` tape image file format is identical to the `.TZX` file format designed by Tomaz Kac.
// Therefore this package is a simple wrapper around the `spectrum/tzx` package.
package cdt

import "github.com/mrcook/retroio/spectrum/tzx"

type CDT struct {
	*tzx.TZX
}
