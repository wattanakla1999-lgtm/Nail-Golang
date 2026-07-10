package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB คือ global variable เก็บ connection ไว้ใช้ทั่วทั้ง project
// ตัวพิมพ์ใหญ่ = Public ใช้ได้จาก package อื่น เช่น handlers
var DB *gorm.DB

func Connect() {
	// DSN = Data Source Name — string บอก GORM ว่าจะต่อ database ที่ไหน อย่างไร
	// format: "host=... user=... password=... dbname=... port=... sslmode=..."
	dsn := "host=localhost user=nailly password=nailly1234 dbname=nailly_db port=5432 sslmode=disable"

	// gorm.Open — เปิด connection กับ database
	// รับ 2 argument: driver (postgres) + config
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// log.Fatal = print error แล้วหยุดโปรแกรมทันที (ต่าง from fmt.Println)
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	fmt.Println("Database connected!")
}
