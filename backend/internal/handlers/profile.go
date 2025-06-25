package handlers

import (
	"brokefolio/backend/internal/middleware"
	"brokefolio/backend/internal/models"
	"brokefolio/backend/internal/route"
	"brokefolio/backend/internal/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ProfileDBHandler struct {
	DB *sql.DB
}

func NewProfileHandler(db *sql.DB) *ProfileDBHandler {
	return &ProfileDBHandler{DB: db}
}

type ChangePasswordReq struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (h *ProfileDBHandler) ProfilePageHandler(w http.ResponseWriter, req *http.Request) {

	isAuthenticated, _ := req.Context().Value(middleware.IsAuthenticatedContextKey).(bool)

	if !isAuthenticated {
		http.Redirect(w, req, "/unauth", http.StatusUnauthorized)
	}

	userID, _ := req.Context().Value(middleware.UserIDContextKey).(string)

	var newUser models.User

	err := h.DB.QueryRow(
		"SELECT name, username, surname, email,pp FROM users WHERE id = $1",
		userID,
	).Scan(&newUser.Name, &newUser.Username, &newUser.Surname, &newUser.Email, &newUser.AvatarURL)

	if err != nil {
		http.Error(w, "SQL Error", http.StatusInternalServerError)
		log.Println("SQL Query error")
		return
	}

	route.RenderTemplate(w, "profile.html", newUser)
}

func (h *ProfileDBHandler) UpdateProfileHandler(w http.ResponseWriter, req *http.Request) {

	isAutheticated, ok := req.Context().Value(middleware.IsAuthenticatedContextKey).(bool)

	if !ok || !isAutheticated {
		http.Error(w, "User Auth Failed", http.StatusUnauthorized)
		log.Println("Auth Context Not Found")
		return
	}

	userID, ok := req.Context().Value(middleware.UserIDContextKey).(string)

	if !ok || userID == "" {
		http.Error(w, "User ID not Found", http.StatusInternalServerError)
		log.Println("User ID not Found")
		return
	}

	// authHeader := req.Header.Get("Authorization")

	// if authHeader == "" {
	// 	http.Error(w, "User Not Authorized", http.StatusUnauthorized)
	// 	log.Println("Unauthorized user!")
	// 	return
	// }

	var updateUser models.User

	err := json.NewDecoder(req.Body).Decode(&updateUser)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		log.Println("Invalid Request On Update")
		return
	}

	if updateUser.Name == "" || updateUser.Surname == "" || updateUser.Username == "" || updateUser.Email == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec("UPDATE users SET name = $1, username = $2, surname = $3, email = $4 WHERE id = $5", updateUser.Name, updateUser.Username, updateUser.Surname, updateUser.Email, userID)

	if err != nil {
		http.Error(w, "Update Failed", http.StatusInternalServerError)
		log.Printf("SQL Update Failed %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func (h *ProfileDBHandler) DeleteProfileHandler(w http.ResponseWriter, req *http.Request) {

	isAutheticated, ok := req.Context().Value(middleware.IsAuthenticatedContextKey).(bool)

	if !ok || !isAutheticated {
		http.Error(w, "User Auth Failed", http.StatusUnauthorized)
		log.Println("Auth Context Not Found")
		return
	}

	userID, ok := req.Context().Value(middleware.UserIDContextKey).(string)

	if !ok || userID == "" {
		http.Error(w, "User ID not Found", http.StatusInternalServerError)
		log.Println("User ID not Found")
		return
	}

	_, err := h.DB.Exec("DELETE FROM users WHERE id = $1", userID)

	if err != nil {
		http.Error(w, "User Could Not Deleted", http.StatusInternalServerError)
		log.Printf("SQL DELETE Error %v\n", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted successfully"})
}

func (h *ProfileDBHandler) ChangePasswordHandler(w http.ResponseWriter, req *http.Request) {

	isAutheticated, ok := req.Context().Value(middleware.IsAuthenticatedContextKey).(bool)

	if !ok || !isAutheticated {
		http.Error(w, "User Auth Failed", http.StatusUnauthorized)
		log.Println("Auth Context Not Found")
		return
	}

	userID, ok := req.Context().Value(middleware.UserIDContextKey).(string)

	if !ok || userID == "" {
		http.Error(w, "User ID not Found", http.StatusInternalServerError)
		log.Println("User ID not Found")
		return
	}

	var cpReq ChangePasswordReq

	err := json.NewDecoder(req.Body).Decode(&cpReq)

	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		log.Printf("Invalid Request %v\n", err)
		return
	}

	if cpReq.CurrentPassword == "" || cpReq.NewPassword == "" {
		http.Error(w, "Mevcut ve yeni şifre boş bırakılamaz.", http.StatusBadRequest)
		return
	}
	if len(cpReq.NewPassword) < 8 {
		http.Error(w, "Yeni şifre en az 8 karakter olmalıdır.", http.StatusBadRequest)
		return
	}

	var hashedPasswordFromDB string
	err = h.DB.QueryRow("SELECT password FROM users WHERE id = $1", userID).Scan(&hashedPasswordFromDB)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Kullanıcı bulunamadı (auth hatası).", http.StatusUnauthorized)
		} else {
			log.Printf("DB error fetching password for user %s: %v", userID, err)
			http.Error(w, "Sunucu hatası: Şifre doğrulanamadı.", http.StatusInternalServerError)
		}
		return
	}

	err = utils.CheckPasswordHash(cpReq.CurrentPassword, hashedPasswordFromDB)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			http.Error(w, "Mevcut şifreniz yanlış.", http.StatusBadRequest)
		} else {
			log.Printf("Bcrypt comparison error for user %s: %v", userID, err)
			http.Error(w, "Sunucu hatası: Şifre doğrulanamadı.", http.StatusInternalServerError)
		}
		return
	}

	newHashedPassword, err := utils.HashPassword(cpReq.NewPassword)

	if err != nil {
		log.Printf("Bcrypt hashing error for user %s: %v", userID, err)
		http.Error(w, "Sunucu hatası: Yeni şifre işlenemedi.", http.StatusInternalServerError)
		return
	}

	_, err = h.DB.Exec("UPDATE users SET password = $1 WHERE id = $2", newHashedPassword, userID)
	if err != nil {
		log.Printf("DB error updating password for user %s: %v", userID, err)
		http.Error(w, "Sunucu hatası: Şifre güncellenemedi.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Şifreniz başarıyla değiştirildi."})

}
