package dto

import (
	"nailly-back-end/internal/model"
	"nailly-back-end/internal/service"
	"time"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token     string        `json:"token"`
	TokenType string        `json:"tokenType"`
	ExpiresAt time.Time     `json:"expiresAt"`
	User      AdminResponse `json:"user"`
}

func ToLoginResponse(result service.AuthResult) LoginResponse {
	return LoginResponse{
		Token: result.Token, TokenType: "Bearer", ExpiresAt: result.ExpiresAt,
		User: ToAdminResponse(result.Admin),
	}
}

func ToAdminResponse(admin model.Admin) AdminResponse {
	return AdminResponse{ID: admin.ID, Username: admin.Username, Name: admin.Name, Role: admin.Role}
}
