package main

import (
	"fmt"
	"time"

	"nailly-back-end/internal/config"
	"nailly-back-end/internal/database"
	"nailly-back-end/internal/model"
	"nailly-back-end/internal/router"
)

func main() {
	setThailandTimezone()

	cfg := config.Load()
	db := database.Connect(cfg.DSN)

	db.AutoMigrate(&model.User{})
	fmt.Println("Database User migrated!")
	db.AutoMigrate(&model.Service{})
	fmt.Println("Database Service migrated!")
	db.AutoMigrate(&model.NailTechnician{})
	fmt.Println("Database Nail Technician migrated!")

	r := router.New(db, cfg.AllowOrigin)
	r.Run(":" + cfg.Port)
}

func setThailandTimezone() {
	time.Local = time.FixedZone("Asia/Bangkok", 7*60*60)
}
