//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var sb strings.Builder
	for len(input) > 0 {
		data, size := utf8.DecodeLastRuneInString(input)
		sb.WriteRune(data)
		input = input[:len(input)-size]
	}
	return sb.String()
}
