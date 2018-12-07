package strings

import (
	"fmt"
	"strings"
)

// Strip0x removes if exists an 0x in front of a string
func Strip0x(s string) string {
	if s[0:2] == "0x" {
		return s[2:]
	}
	return s
}

// Add0x adds if not exists an 0x in front of a string
func Add0x(s string) string {
	if s[0:2] != "0x" {
		return "0x" + s
	}
	return s
}

// https://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go
func BuildString(s ...string) string {

	var b strings.Builder
	for _, v := range s {
		fmt.Fprintf(&b, v)
	}

	return b.String()
}
