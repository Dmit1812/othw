package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidString    = errors.New("invalid string")
	ErrUnsupportedDigit = errors.New("unsupported digit in rune to uint conversion")
)

type runeType uint8

const (
	escapecharacter runeType = iota
	digit
	other
)

const (
	backslash     rune = '\\'
	zero          rune = '0'
	zerofullwidth rune = 0xFF10
)

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

func runeDigitToInt(r rune) (int, error) {
	if r >= 0x0030 && r <= 0x0039 {
		return int(r - zero), nil
	}
	if r >= 0xFF10 && r <= 0xFF19 {
		return int(r - zerofullwidth), nil
	}
	return 0, ErrUnsupportedDigit
}

func runeDetermineType(r rune) runeType {
	if isDigit(r) {
		return digit
	}
	if r == backslash {
		return escapecharacter
	}
	return other
}

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
				return "", ErrInvalidString
			}
			// write rune according to rules of digits
			// получаем число из руны
			count, err := runeDigitToInt(r)
			// сигнализируем об ошибке если была ошибка в переводе цифры в число
			if err != nil {
				return "", ErrInvalidString
			}

			// повторяем столько сколько нужно
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

			// if this is the last rune from string error as we need another for escaping
			if len(str) == size {
				return "", ErrInvalidString
			}

			// then read next character and place it in the waiting
			str = str[size:]
			r, size = utf8.DecodeRuneInString(str)

			// if next character class is not a number or slash error
			if runeDetermineType(r) == other {
				return "", ErrInvalidString
			}

			// place the new character in the waiting
			runeToWrite = r
			runeWaiting = true
		case other:
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
