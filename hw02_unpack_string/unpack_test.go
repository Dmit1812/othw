package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "ф5a赓3c", expected: "фффффa赓赓赓c"},
		{input: "a1cc", expected: "acc"},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "qw\\ne", "qw\\"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		input    rune
		expected bool
	}{
		{input: '0', expected: true},
		{input: 'a', expected: false},
		{input: '9', expected: true},
		{input: 0xFF10, expected: true},
		{input: 0xFF19, expected: true},
		{input: 0xFF20, expected: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(string(tc.input), func(t *testing.T) {
			result := isDigit(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestRuneDigitToInt(t *testing.T) {
	tests := []struct {
		input       rune
		expected    int
		shouldError bool
	}{
		{input: '0', expected: 0, shouldError: false},
		{input: 'a', expected: 0, shouldError: true},
		{input: '9', expected: 9, shouldError: false},
		{input: 0xFF10, expected: 0, shouldError: false},
		{input: 0xFF19, expected: 9, shouldError: false},
		{input: 0xFF20, expected: 0, shouldError: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(string(tc.input), func(t *testing.T) {
			result, err := runeDigitToInt(tc.input)
			if tc.shouldError {
				require.Truef(t, errors.Is(err, ErrUnsupportedDigit), "actual error %q", err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestRuneDetermineType(t *testing.T) {
	tests := []struct {
		input    rune
		expected runeType
	}{
		{input: '0', expected: digit},
		{input: 'a', expected: other},
		{input: '9', expected: digit},
		{input: '\\', expected: escapecharacter},
		{input: 0xFF10, expected: digit},
		{input: 0xFF19, expected: digit},
		{input: 0xFF20, expected: other},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(string(tc.input), func(t *testing.T) {
			result := runeDetermineType(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
