package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var stringBuilder strings.Builder
	var repeatCount int
	var stringLength int

	runes := []rune(s)
	stringLength = len(runes)

	for i := 0; i < stringLength; i++ {
		if unicode.IsNumber(runes[i]) {
			return "", ErrInvalidString
		}

		if runes[i] == '\\' {
			i++
		}

		if stringLength > i+1 && unicode.IsNumber(runes[i+1]) {
			repeatCount, _ = strconv.Atoi(string(runes[i+1]))
			stringBuilder.WriteString(strings.Repeat(string(runes[i]), repeatCount))
			i++
		} else {
			stringBuilder.WriteString(string(runes[i]))
		}
	}

	return stringBuilder.String(), nil
}
