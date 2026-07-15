package service

import (
	"errors"
	"fmt"
	"nailly-back-end/internal/apperror"
	"nailly-back-end/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const jwtIssuer = "nailly-api"

type AuthStore interface {
	FindAdminByUsername(username string) (model.Admin, error)
	CreateAdmin(admin *model.Admin) error
	SaveAdmin(admin *model.Admin) error
}

type AdminClaims struct {
	AdminID  uint   `json:"adminId"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl}
}

func (m *JWTManager) Generate(admin model.Admin) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)
	claims := AdminClaims{
		AdminID: admin.ID, Username: admin.Username, Name: admin.Name, Role: admin.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: jwtIssuer, Subject: strconv.FormatUint(uint64(admin.ID), 10),
			IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	return token, expiresAt, err
}

func (m *JWTManager) Verify(tokenValue string) (*AdminClaims, error) {
	claims := &AdminClaims{}
	token, err := jwt.ParseWithClaims(
		tokenValue,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
			}
			return m.secret, nil
		},
		jwt.WithIssuer(jwtIssuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	return claims, nil
}

type AuthResult struct {
	Admin     model.Admin
	Token     string
	ExpiresAt time.Time
}

type AuthService struct {
	repo       AuthStore
	jwtManager *JWTManager
}

func NewAuthService(repo AuthStore, jwtManager *JWTManager) *AuthService {
	return &AuthService{repo: repo, jwtManager: jwtManager}
}

func (s *AuthService) Login(username, password string) (AuthResult, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	if username == "" || password == "" {
		return AuthResult{}, apperror.BadRequest("username and password are required", apperror.ErrValidation)
	}
	admin, err := s.repo.FindAdminByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) || err == nil && bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return AuthResult{}, apperror.Unauthorized("ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง", apperror.ErrValidation)
	}
	if err != nil {
		return AuthResult{}, err
	}
	if !admin.Active {
		return AuthResult{}, apperror.Unauthorized("บัญชีผู้ดูแลถูกปิดใช้งาน", apperror.ErrValidation)
	}
	token, expiresAt, err := s.jwtManager.Generate(admin)
	if err != nil {
		return AuthResult{}, apperror.Internal("could not create access token", err)
	}
	return AuthResult{Admin: admin, Token: token, ExpiresAt: expiresAt}, nil
}

func (s *AuthService) EnsureAdmin(username, name, password string) error {
	username = strings.ToLower(strings.TrimSpace(username))
	if username == "" || strings.TrimSpace(name) == "" || password == "" {
		return errors.New("configured admin username, name and password are required")
	}

	admin, err := s.repo.FindAdminByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		return s.repo.CreateAdmin(&model.Admin{
			Username: username, Name: strings.TrimSpace(name), PasswordHash: string(hash), Role: "admin", Active: true,
		})
	}
	if err != nil {
		return err
	}

	changed := admin.Name != strings.TrimSpace(name) || admin.Role != "admin" || !admin.Active
	admin.Name, admin.Role, admin.Active = strings.TrimSpace(name), "admin", true
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		admin.PasswordHash = string(hash)
		changed = true
	}
	if changed {
		return s.repo.SaveAdmin(&admin)
	}
	return nil
}
