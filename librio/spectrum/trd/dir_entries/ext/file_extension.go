package ext

// File Extension is a character that describes the file type.
type FileExtension byte

const (
	Basic      FileExtension = 'B' // B: BASIC Program
	Data       FileExtension = 'D' // D: DATA array (numeric or alphanumeric)
	Code       FileExtension = 'C' // C: CODE
	AccessType FileExtension = '#' // #: Print file (may be split into several sub-files with max. 4096 bytes each)
)
