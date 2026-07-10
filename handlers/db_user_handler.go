package handlers

import (
	"encoding/json"
	"fmt"
	"nailly-back-end/database"
	"nailly-back-end/models"
	"nailly-back-end/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/users
func DBGetUsers(c *gin.Context) {

	var users []models.UserDB
	query := database.DB.Model(&models.UserDB{})
	query = utils.ApplyLikeFilters(c, query, map[string]string{
		"name":  "name",
		"email": "email",
	})

	pagination, total, err := utils.Paginate(c, query, &users)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"page":  pagination.Page,
		"limit": pagination.Limit,
		"total": total,
	})
}

// GET /api/users/:id
func DBGetUserByID(c *gin.Context) {
	var user models.UserDB
	// DB.First(&user, id) = SELECT * FROM user_dbs WHERE id = ? LIMIT 1
	// .Error — GORM return error ถ้าหาไม่เจอ (record not found)
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func DBGetUserByEmail(c *gin.Context) {
	var user models.UserDB
	email := c.Param("email")
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func DBGetUsersOlderThan(c *gin.Context) {
	var users []models.UserDB
	age, err := strconv.Atoi(c.Param("age"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid age"})
		return
	}

	if err := database.DB.Where("age > ?", age).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// POST /api/users
func DBCreateUser(c *gin.Context) {
	var input models.UserDB
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if input.Age <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "age is required and must be greater than 0"})
		return
	}

	bodyPretty, _ := json.MarshalIndent(input, "", "  ")
	fmt.Println("Request Body:", string(bodyPretty))

	// DB.Create(&input) = INSERT INTO user_dbs (...) VALUES (...)
	// GORM จะเติม ID, CreatedAt, UpdatedAt ให้อัตโนมัติ
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// PUT /api/users/:id
func DBUpdateUser(c *gin.Context) {
	var user models.UserDB
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var input models.UserDB
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// DB.Model(&user).Updates(input) = UPDATE user_dbs SET ... WHERE id = ?
	// .Model() บอกว่า update record ไหน
	// .Updates() อัปเดตเฉพาะ field ที่มีค่า (ไม่ใช่ทุก field)
	if err := database.DB.Model(&user).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DELETE /api/users/:id
func DBDeleteUser(c *gin.Context) {
	var user models.UserDB
	if err := database.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	// DB.Delete(&user) = Soft Delete — ไม่ได้ลบจริงๆ
	// แค่ตั้งค่า deleted_at = now() ใน DB
	// ถ้าอยาก hard delete ต้องใช้ DB.Unscoped().Delete(&user)

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
