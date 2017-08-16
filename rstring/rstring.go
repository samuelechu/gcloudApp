package rstring

import "strings"

func Reverse(s string) string{
     b := []rune(s)
     for i := 0; i < len(b)/2; i++{
     	 j:= len(b)-i-1
	 b[i], b[j] = b[j], b[i]
     }

     return string(b)
}

func RemoveWhitespace(s string) string{
	fields := strings.Fields(s)

	result := ""

	for _, value := range fields {
		result += value
	}

	return result
}