package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidString    = errors.New("invalid string")
	ErrUnsupportedDigit = errors.New("unsupported digit in rune to uint conversion")
)

type operation uint8

const (
	getRune operation = iota
	getRepeatCountOrWrite
	getDigitOrSlash
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

func Unpack(str string) (string, error) {
	var op operation
	var result strings.Builder
	var runeToWrite rune
	var runeWaiting bool
	fmt.Printf("%s\n", str)
	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		// fmt.Printf("string: %s, runeS: %s, runeD:%d, size: %d, op: %d, runeToWrite: %s, runeWaiting: %t\n",
		//             result.String(), string(r), r, size, op, string(runeToWrite), runeWaiting)
		switch op {
		case getRune:
			if isDigit(r) {
				// если это цифра - ошибка
				return "", ErrInvalidString
			}

			// Если есть руна на запись пишем
			if runeWaiting {
				result.WriteRune(runeToWrite)
			}

			// если слеш, то включаем getDigitOrSlash
			if r == backslash {
				op = getDigitOrSlash
			} else {
				runeToWrite = r
				runeWaiting = true
				op = getRepeatCountOrWrite
			}
		case getRepeatCountOrWrite:
			// если не цифра то записываем
			if !isDigit(r) {
				// если есть руна на запись, пишем
				if runeWaiting {
					result.WriteRune(runeToWrite)
					runeToWrite = r
					runeWaiting = true
					op = getRepeatCountOrWrite
				} else {
					return "", ErrInvalidString
				}
				if r == backslash {
					// если слеш, то включаем getDigitOrSlash
					op = getDigitOrSlash
				}
			} else {
				// получаем число из руны
				count, err := runeDigitToInt(r)
				// сигнализируем об ошибке если была ошибка в переводе цифры в число
				if err != nil {
					return "", ErrInvalidString
				}

				// Если число повторений 0 не будем добавлять руну
				if count == 0 {
					runeWaiting = false
				} else {
					// иначе повторяем столько сколько нужно
					for i := 0; i < count; i++ {
						result.WriteRune(runeToWrite)
					}
					runeWaiting = false
				}
				op = getRune
			}
		case getDigitOrSlash:
			// если цифра или слеш то добавляем в результат
			// иначе ошибка
			if isDigit(r) || r == backslash {
				runeToWrite = r
				runeWaiting = true
			} else {
				return "", ErrInvalidString
			}
			op = getRepeatCountOrWrite
		}
		str = str[size:]
	}
	// write the last rune if present
	if runeWaiting {
		result.WriteRune(runeToWrite)
	}
	return result.String(), nil
}
