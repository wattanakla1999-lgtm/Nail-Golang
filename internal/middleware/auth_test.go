package middleware

import (
	"nailly-back-end/internal/model"
	"nailly-back-end/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestRequireAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	manager := service.NewJWTManager("middleware-test-secret", time.Hour)
	token, _, err := manager.Generate(model.Admin{
		Model: gorm.Model{ID: 7}, Username: "admin", Name: "Admin", Role: "admin",
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	router := gin.New()
	router.GET("/protected", RequireAdmin(manager), func(c *gin.Context) {
		claims, ok := AdminClaimsFromContext(c)
		if !ok || claims.AdminID != 7 {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})

	withoutToken := httptest.NewRecorder()
	router.ServeHTTP(withoutToken, httptest.NewRequest(http.MethodGet, "/protected", nil))
	if withoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", withoutToken.Code)
	}

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	withToken := httptest.NewRecorder()
	router.ServeHTTP(withToken, request)
	if withToken.Code != http.StatusOK {
		t.Fatalf("with token status = %d, want 200", withToken.Code)
	}
}
