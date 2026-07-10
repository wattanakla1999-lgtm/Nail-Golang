package main

import (
	"fmt"
	"nailly-back-end/database"
	"nailly-back-end/handlers"
	"nailly-back-end/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// เชื่อมต่อ PostgreSQL ผ่าน GORM
	database.Connect()

	// AutoMigrate — GORM อ่าน struct แล้วสร้าง/อัปเดต table ใน DB อัตโนมัติ
	// ถ้า table ยังไม่มี → สร้างใหม่
	// ถ้า table มีแล้ว → เพิ่ม column ที่ขาด (ไม่ลบ column เก่า)
	database.DB.AutoMigrate(&models.UserDB{})
	fmt.Println("Database migrated!")

	r := gin.Default()

	// Routes in-memory (บทเรียนก่อนหน้า)
	r.GET("/users", handlers.GinGetUsers)
	r.POST("/users", handlers.GinCreateUser)
	r.GET("/users/:id", handlers.GinGetUserByID)
	r.PUT("/users/:id", handlers.GinUpdateUser)
	r.DELETE("/users/:id", handlers.GinDeleteUser)

	// Route Group — จัดกลุ่ม routes ที่มี prefix เดียวกัน
	// r.Group("/api") = ทุก route ข้างในจะขึ้นต้นด้วย /api
	// เช่น api.GET("/users") → /api/users
	api := r.Group("/api")
	api.GET("/users", handlers.DBGetUsers)
	api.GET("/users/email/:email", handlers.DBGetUserByEmail)
	api.GET("/users/age/:age", handlers.DBGetUsersOlderThan)
	api.GET("/users/:id", handlers.DBGetUserByID)

	api.POST("/users", handlers.DBCreateUser)
	api.PUT("/users/:id", handlers.DBUpdateUser)
	api.DELETE("/users/:id", handlers.DBDeleteUser)

	r.Run(":8080")
}
