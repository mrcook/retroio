package trd

import "fmt"

type FileType interface {
	Name() string
	Info(info FileInformation) string
}

type FileTypeBasic struct{}

func (t *FileTypeBasic) Name() string {
	return "BASIC Program"
}

func (t *FileTypeBasic) Info(i FileInformation) string {
	return ""
}

type FileTypeCode struct{}

func (t *FileTypeCode) Name() string {
	return "Code (bytes)"
}

func (t *FileTypeCode) Info(i FileInformation) string {
	return fmt.Sprintf(" - Address:     %d\n", i.StartAddress)
}

type FileTypeOther struct{
	Extension byte
}

func (t *FileTypeOther) Name() string {
	return fmt.Sprintf("Other type (%c)", t.Extension)
}

func (t *FileTypeOther) Info(i FileInformation) string {
	return ""
}
