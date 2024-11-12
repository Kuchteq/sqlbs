package main

import (
	"regexp"
	"strings"
)

func ApoQuote(s string) string {
	return "'" + s + "'"
}
func UnparenthesizeAndTrim(s string) string {
	if len(s) >= 2 && strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = s[1 : len(s)-1]
		s = strings.TrimSpace(s)
	}
	return s
}
func WithinParenthesis(s string) string {
	re := regexp.MustCompile(`(?s)\((\S.*)\)`)
	match := re.FindStringSubmatch(s)
	if len(match) < 1 {
		return ""
	}
	return match[1]

}
func findFieldByName(fields []*Field, name string) *Field {
	for _, field := range fields {
		if field.Name == name {
			return field
		}
	}
	return nil // Return nil if not found
}
