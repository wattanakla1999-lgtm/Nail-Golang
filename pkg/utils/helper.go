package utils

import "strings"

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func trimSpaces(s string) string {
	return strings.TrimSpace(s)
}

func SanitizeName(name string) string {
	return Capitalize(trimSpaces(name))
}
