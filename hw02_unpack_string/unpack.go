package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/runenames"
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
	return unicode.IsDigit(r)
}

// runeDigitToInt converts a rune if it is of supported types to integer
// supported types are any unicode.IsDigit()
// получается что ни strconv.Atoi(), ни strconv.ParseInt()
// не поддерживают все цифры из алфавита utf-8
// вот единственный способ, что удалось найти :)
func runeDigitToInt(r rune) (int, error) {
	var name = runenames.Name(r)
	switch {
	case strings.Contains(name, "DIGIT ZERO"):
		return 0, nil
	case strings.Contains(name, "DIGIT ONE"):
		return 1, nil
	case strings.Contains(name, "DIGIT TWO"):
		return 2, nil
	case strings.Contains(name, "DIGIT THREE"):
		return 3, nil
	case strings.Contains(name, "DIGIT FOUR"):
		return 4, nil
	case strings.Contains(name, "DIGIT FIVE"):
		return 5, nil
	case strings.Contains(name, "DIGIT SIX"):
		return 6, nil
	case strings.Contains(name, "DIGIT SEVEN"):
		return 7, nil
	case strings.Contains(name, "DIGIT EIGHT"):
		return 8, nil
	case strings.Contains(name, "DIGIT NINE"):
		return 9, nil
	default:
		return 0, ErrUnsupportedDigit
	}
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
