package models

import "gorm.io/gorm"

type UserDB struct {
	// gorm.Model = Embedded struct จาก GORM
	// เพิ่ม field ให้อัตโนมัติ:
	//   ID        uint      — primary key, auto increment
	//   CreatedAt time.Time — วันที่สร้าง
	//   UpdatedAt time.Time — วันที่แก้ไขล่าสุด
	//   DeletedAt *time.Time — Soft Delete (ไม่ลบจริง แค่ mark)
	gorm.Model

	// Struct Tags — metadata ที่ติดกับ field ใช้ backtick ครอบ
	// `json:"name"`     → ตอน encode/decode JSON ใช้ชื่อ "name" (ตัวเล็ก)
	// `gorm:"not null"` → สร้าง column ใน DB แบบ NOT NULL constraint
	Name string `json:"name" gorm:"not null"`

	// gorm:"uniqueIndex" → สร้าง UNIQUE INDEX ใน DB (email ซ้ำไม่ได้)
	Email string `json:"email" gorm:"uniqueIndex;not null"`

	Age int `json:"age"`
}
