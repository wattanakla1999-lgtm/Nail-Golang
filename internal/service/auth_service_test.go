package service

import (
	"errors"
	"nailly-back-end/internal/apperror"
	"nailly-back-end/internal/model"
	"net/http"
	"testing"
	"time"

	"gorm.io/gorm"
)

type fakeAuthStore struct {
	admins map[string]model.Admin
	nextID uint
}

func newFakeAuthStore() *fakeAuthStore {
	return &fakeAuthStore{admins: make(map[string]model.Admin), nextID: 1}
}

func (f *fakeAuthStore) FindAdminByUsername(username string) (model.Admin, error) {
	admin, ok := f.admins[username]
	if !ok {
		return model.Admin{}, gorm.ErrRecordNotFound
	}
	return admin, nil
}

func (f *fakeAuthStore) CreateAdmin(admin *model.Admin) error {
	admin.ID = f.nextID
	f.nextID++
	f.admins[admin.Username] = *admin
	return nil
}

func (f *fakeAuthStore) SaveAdmin(admin *model.Admin) error {
	f.admins[admin.Username] = *admin
	return nil
}

func TestAuthServiceLoginAndVerifyToken(t *testing.T) {
	store := newFakeAuthStore()
	manager := NewJWTManager("test-secret-that-is-long-enough", time.Hour)
	authService := NewAuthService(store, manager)
	if err := authService.EnsureAdmin("admin", "ผู้ดูแลระบบ", "nailly2025"); err != nil {
		t.Fatalf("EnsureAdmin() error = %v", err)
	}

	result, err := authService.Login(" ADMIN ", "nailly2025")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	claims, err := manager.Verify(result.Token)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if claims.AdminID != result.Admin.ID || claims.Username != "admin" || claims.Role != "admin" {
		t.Fatalf("claims = %+v", claims)
	}
}

func TestAuthServiceRejectsInvalidCredentials(t *testing.T) {
	store := newFakeAuthStore()
	authService := NewAuthService(store, NewJWTManager("test-secret", time.Hour))
	if err := authService.EnsureAdmin("admin", "Admin", "correct-password"); err != nil {
		t.Fatalf("EnsureAdmin() error = %v", err)
	}

	_, err := authService.Login("admin", "wrong-password")
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Status != http.StatusUnauthorized {
		t.Fatalf("Login() error = %v, want 401 AppError", err)
	}
}

func TestJWTManagerRejectsExpiredToken(t *testing.T) {
	manager := NewJWTManager("test-secret", -time.Minute)
	token, _, err := manager.Generate(model.Admin{Model: gorm.Model{ID: 1}, Username: "admin", Role: "admin"})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if _, err := manager.Verify(token); err == nil {
		t.Fatal("Verify() error = nil, want expired token error")
	}
}
