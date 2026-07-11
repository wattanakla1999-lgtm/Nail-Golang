package main

import (
	"fmt"
	"nailly-back-end/internal/config"
	"nailly-back-end/internal/database"
	"nailly-back-end/internal/model"
	"nailly-back-end/internal/router"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DSN)

	db.AutoMigrate(&model.User{})
	fmt.Println("Database migrated!")

	r := router.New(db)
	r.Run(":" + cfg.Port)
}
