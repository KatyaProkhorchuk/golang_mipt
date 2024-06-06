//go:build !solution

package varfmt

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Sprintf(format string, args ...interface{}) string {
	length := 0
	argsString := make([]string, length, 10) // len, cap
	var sb strings.Builder
	for _, arg := range args {
		str := fmt.Sprint(arg)
		length += len(str)
		argsString = append(argsString, str)
	}
	sb.Grow(length + len(format)) // 15 alloc -> 4 alloc for big test :)
	num := 0
	ind := 0
	data := -1
	for len(format) > 0 {
		r, size := utf8.DecodeRuneInString(format)
		format = format[size:] // уменьшаем format
		if r == '{' {
			data = 0
		} else if r == '}' {
			if data == 0 {
				ind = num
			}
			// fmt.Println(data)
			sb.WriteString(argsString[ind])
			num++
			data = -1
			ind = 0
		} else if data >= 0 {
			ind = ind*10 + (int(r) - '0')
			data++
		} else {
			sb.WriteRune(r)
		}
	}
	fmt.Println(sb.String())
	return sb.String()
}
