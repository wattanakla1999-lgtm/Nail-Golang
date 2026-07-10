package handlers

import (
	"encoding/json"
	"fmt"
	"nailly-back-end/models"
	"net/http"
	"strconv"
	"strings"
)

// จำลอง in-memory database
var users = []*models.User{
	models.NewUser(1, "Kla", "kla@example.com", 28),
	models.NewUser(2, "Bob", "bob@example.com", 25),
	models.NewUser(3, "Alice", "alice@example.com", 30),
}

// writeJSON — helper ส่ง JSON response
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// GET /users — ดึง user ทั้งหมด
func GetUsers(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, users)
}

// GET /users/{id} — ดึง user ตาม ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	// ตัด /users/ ออก เหลือแค่ id
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	for _, u := range users {
		if u.ID == id {
			writeJSON(w, http.StatusOK, u)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
}

// POST /users — สร้าง user ใหม่
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var input models.User

	fmt.Println("Request Body:", r.Body) // Debug: print the request body

	fmt.Println("Request Method:", r.Method) // Debug: print the request method

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	input.ID = len(users) + 1
	users = append(users, &input)
	writeJSON(w, http.StatusCreated, &input)
}

// PUT /users/{id} — อัปเดต user ตาม ID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// ตัด /users/ ออก เหลือแค่ id
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	var input models.User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	for _, u := range users {
		if u.ID == id {
			u.Name = input.Name
			u.Email = input.Email
			u.Age = input.Age
			writeJSON(w, http.StatusOK, u)
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// ตัด /users/ ออก เหลือแค่ id
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
			return
		}
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
}

// UserRouter — จัดการ routing ของ /users และ /users/{id}
func UserRouter(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// /users หรือ /users/
	if path == "/users" || path == "/users/" {
		switch r.Method {
		case http.MethodGet:
			GetUsers(w, r)
		case http.MethodPost:
			CreateUser(w, r)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
		return
	}

	// /users/{id}
	if strings.HasPrefix(path, "/users/") {
		switch r.Method {
		case http.MethodGet:
			GetUserByID(w, r)
		case http.MethodPut:
			UpdateUser(w, r)
		case http.MethodDelete:
			DeleteUser(w, r)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
		return
	}

	writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}
