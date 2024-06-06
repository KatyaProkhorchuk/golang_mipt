//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var sb strings.Builder
	isPrevSpase := false
	for len(input) > 0 {
		data, size := utf8.DecodeRuneInString(input)
		isSpase := unicode.IsSpace(data)
		if !isSpase {
			sb.WriteRune(data)
		} else if !isPrevSpase {
			sb.WriteRune(' ')
		}
		isPrevSpase = isSpase
		input = input[size:]
	}
	return sb.String()
}
