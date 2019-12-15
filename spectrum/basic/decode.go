// Package basic is a simple decoder for ZX Spectrum BASIC as stored on tapes.
package basic

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Decode as ZX Spectrum BASIC program
func Decode(programData []byte) ([]string, error) {
	var basic []string

	reader := bytes.NewReader(programData)
	for {
		lineNum, lineLen, err := getLineNumAndLen(reader)

		// Check if we've read all the data
		if err != nil {
			if err == io.EOF {
				// all done, return the results
				return basic, nil
			} else if err != io.ErrUnexpectedEOF {
				return nil, err
			}
		}

		data := make([]byte, lineLen)
		_, err = reader.Read(data)
		if err != nil && err != io.EOF {
			return nil, err
		}

		basicString := decodeBasicBytes(data)
		line := fmt.Sprintf("%4d %s", lineNum, basicString)

		basic = append(basic, line)
	}
}

func getLineNumAndLen(reader *bytes.Reader) (lineNumber, lineLength uint16, err error) {
	tmpBytes := []byte{0, 0}

	// If we don't have an error here then we're safe to read the next
	// two bytes without checking for errors.
	if _, err = reader.Read(tmpBytes); err != nil {
		return
	}
	lineNumber = BigEndianToInt(tmpBytes)

	_, _ = reader.Read(tmpBytes)
	lineLength = LittleEndianToInt(tmpBytes)

	return
}

// Decodes a line of bytes from a BASIC program
func decodeBasicBytes(lineOfBasic []byte) string {
	pos := 0
	length := len(lineOfBasic)

	basic := ""
	var lastCharOfLine byte

	for pos = 0; pos < length; {
		if len(basic) > 0 {
			lastCharOfLine = basic[len(basic)-1]
		}
		char := lineOfBasic[pos]
		pos += 1

		switch {
		case char >= 0x10 && char <= 0x15:
			pos += 1
		case char == 0x16, char == 0x17:
			pos += 2
		case char == 0x0E:
			pos += 5
		default:
			basic += decodeWithPadding(char, lastCharOfLine)
		}
	}

	return basic
}

// TODO: this needs improving, but is functional for the moment
func decodeWithPadding(char, lastChar byte) string {
	decoded := CharacterSet[char]

	// No padding
	if char < 0xA5 {
		return decoded
	}

	if lastChar != ' ' {
		// SCREEN$ and CODE have no space before them if they follow a double quote.
		if lastChar != '"' && char != 0xAA && char != 0xAF {
			decoded = " " + decoded
		}
	}

	return decoded + " "
}

// BigEndianToInt converts a slice of 2 big endian ordered bytes.
func BigEndianToInt(b []byte) uint16 {
	return binary.BigEndian.Uint16(b[:])
}

// LittleEndianToInt converts a slice of 2 little endian ordered bytes.
func LittleEndianToInt(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b[:])
}
