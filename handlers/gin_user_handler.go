package handlers

import (
	"nailly-back-end/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// จำลอง in-memory database
var ginUsers = []*models.User{
	models.NewUser(1, "Kla", "kla@example.com", 28),
	models.NewUser(2, "Bob", "bob@example.com", 25),
	models.NewUser(3, "Alice", "alice@example.com", 30),
}

// GET /users
func GinGetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, ginUsers)
}

// GET /users/:id
func GinGetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) // ดึง :id จาก URL ง่ายมาก!
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	for _, u := range ginUsers {
		if u.ID == id {
			c.JSON(http.StatusOK, u)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}

// POST /users
func GinCreateUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil { // validate + decode อัตโนมัติ
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = len(ginUsers) + 1
	ginUsers = append(ginUsers, &input)
	c.JSON(http.StatusCreated, &input)
}

// PUT /users/:id
func GinUpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, u := range ginUsers {
		if u.ID == id {
			u.Name = input.Name
			u.Email = input.Email
			u.Age = input.Age
			c.JSON(http.StatusOK, u)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}

// DELETE /users/:id
func GinDeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	for i, u := range ginUsers {
		if u.ID == id {
			ginUsers = append(ginUsers[:i], ginUsers[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}
