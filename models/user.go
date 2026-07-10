package models

import "fmt"

// User struct อยู่ใน package "models"
type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

// Constructor function — สร้าง User ใหม่ (Go ไม่มี new keyword แบบ Java)
func NewUser(id int, name, email string, age int) *User {
	return &User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}
}

func (u *User) String() string {
	return fmt.Sprintf("[%d] %s <%s> age:%d", u.ID, u.Name, u.Email, u.Age)
}

func (u *User) IsAdult() bool {
	return u.Age >= 18
}
