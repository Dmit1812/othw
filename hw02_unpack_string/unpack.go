package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrUnsupportedDigit = errors.New("a non supported digit or rune was used to convert to integer")
	ErrBackslashUsage   = errors.New("there should be a backslash or a digit after a backslash")
	ErrDigitNotInPlace  = errors.New("digits can't be used as first character or after another digit, " +
		"consider escaping with backslash")
)

// define runeType and values we distinguish.
type runeType uint8

const (
	escapecharacter runeType = iota
	digit
	other
)

// define runes.
const (
	backslash     rune = '\\'
	zero          rune = '0'
	zerofullwidth rune = 0xFF10
)

// isDigit determines if a rune is a supported digit.
func isDigit(r rune) bool {
	// we consider as digits only normal and fullwidth digits
	if r >= 0x0030 && r <= 0x0039 {
		return true
	}
	if r >= 0xFF10 && r <= 0xFF19 {
		return true
	}
	return false
}

// runeDigitToInt converts a rune if it is of supported types to integer
// supported types are normal and fullwidth digits.
func runeDigitToInt(r rune) (int, error) {
	if r >= 0x0030 && r <= 0x0039 {
		return int(r - zero), nil
	}
	if r >= 0xFF10 && r <= 0xFF19 {
		return int(r - zerofullwidth), nil
	}
	return 0, ErrUnsupportedDigit
}

// runeDetermineType return which type rune belongs to digit, escapecharacter or other.
func runeDetermineType(r rune) runeType {
	if isDigit(r) {
		return digit
	}
	if r == backslash {
		return escapecharacter
	}
	return other
}

// Unpack decodes a string according to requirements
// specified here: https://github.com/OtusGolang/home_work/blob/master/hw02_unpack_string/README.md .
func Unpack(str string) (string, error) {
	var result strings.Builder
	var runeToWrite rune
	var runeWaiting bool

	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		runeClass := runeDetermineType(r)

		switch runeClass {
		case digit:
			// if we have no previous rune and it is a digit - error
			if !runeWaiting {
				return "", ErrDigitNotInPlace
			}

			// get a count of times to replicate a rune in result string
			count, err := runeDigitToInt(r)
			// report error in digit conversion if there was one
			if err != nil {
				return "", ErrUnsupportedDigit
			}

			// replicate rune needed number of times
			for i := 0; i < count; i++ {
				result.WriteRune(runeToWrite)
			}
			runeWaiting = false
		case escapecharacter:
			// write previous character if it is there
			if runeWaiting {
				// write previous character if it is there and remember next one
				result.WriteRune(runeToWrite)
			}

			// if this is the last rune from string error as we need another character for escaping
			if len(str) == size {
				return "", ErrBackslashUsage
			}

			// then read next rune and determine it's type
			str = str[size:]
			r, size = utf8.DecodeRuneInString(str)

			// if next rune class is not a number or slash - error
			if runeDetermineType(r) == other {
				return "", ErrBackslashUsage
			}

			// place the new character in runeToWrite
			runeToWrite = r
			runeWaiting = true
		case other:
			// if we have unwritten rune write it
			if runeWaiting {
				// write previous character if it is there and remember next one
				result.WriteRune(runeToWrite)
			}
			runeToWrite = r
			runeWaiting = true
		}
		str = str[size:]
	}

	// write the last rune if present
	if runeWaiting {
		result.WriteRune(runeToWrite)
	}

	return result.String(), nil
}
