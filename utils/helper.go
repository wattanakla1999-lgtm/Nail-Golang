package utils

import "strings"

// ตัวพิมพ์ใหญ่ = public (ใช้ได้จาก package อื่น)
// ตัวพิมพ์เล็ก = private (ใช้ได้แค่ใน package นี้)

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ฟังก์ชัน private — ใช้ได้แค่ใน package utils เท่านั้น
func trimSpaces(s string) string {
	return strings.TrimSpace(s)
}

func SanitizeName(name string) string {
	return Capitalize(trimSpaces(name)) // เรียกใช้ private func ภายใน package ได้
}
