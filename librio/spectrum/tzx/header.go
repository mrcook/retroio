package tzx

import "github.com/pkg/errors"

// New TZX should always default to latest version.
const (
	supportedMajorVersion = 1
	supportedMinorVersion = 20
)

// TZX Header
//
// The header is the first block of data found in all TZX files. The file is
// identified with the first 7 bytes being `ZXTape!`, followed by the _end
// of file_ byte `26` (`1A` hex). This is then followed by two bytes containing
// the major and minor version numbers of the TZX specification used.
type Header struct {
	Signature    [8]byte // must be `ZXTape!` + 0x1A terminator byte
	MajorVersion uint8   // TZX major revision number
	MinorVersion uint8   // TZX minor revision number
}

// Validates the TZX header data.
func (h Header) valid() error {
	var validationError error

	if string(h.Signature[:6]) != "ZXTape!" {
		validationError = errors.Wrapf(validationError, "Incorrect signature, got '%s'", h.Signature)
	}

	if h.Signature[7] != 0x1a {
		validationError = errors.Wrapf(validationError, "Incorrect terminator, got '%b'", h.Signature[7])
	}

	if h.MajorVersion != supportedMajorVersion {
		validationError = errors.Wrapf(validationError, "Invalid version, got v%d.%d", h.MajorVersion, h.MinorVersion)
	}

	return validationError
}
