package handlers

import (
	"brokefolio/backend/internal/middleware"
	"brokefolio/backend/internal/models"
	"brokefolio/backend/internal/route"
	"brokefolio/backend/internal/utils"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
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
		http.Redirect(w, req, "/unauthorized", http.StatusUnauthorized)
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

	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Kayıt verisi işlenemedi.", http.StatusBadRequest)
		return
	}

	name := req.FormValue("name")
	surname := req.FormValue("surname")
	username := req.FormValue("username")
	email := req.FormValue("email")

	if name == "" || surname == "" || username == "" || email == "" {
		utils.SendJSONError(w, "Tüm alanlar doldurulmalıdır.", http.StatusBadRequest)
		return
	}

	var existingCount int
	err = h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE (username = $1 OR email = $2) AND id != $3", username, email, userID).Scan(&existingCount)
	if err != nil {
		log.Printf("Error checking existing username/email for update: %v", err)
		utils.SendJSONError(w, "Sunucu hatası. Lütfen tekrar deneyin.", http.StatusInternalServerError)
		return
	}
	if existingCount > 0 {
		var usernameExists, emailExists bool
		h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1 AND id != $2", username, userID).Scan(&usernameExists)
		h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1 AND id != $2", email, userID).Scan(&emailExists)

		if usernameExists {
			utils.SendJSONError(w, "Bu kullanıcı adı zaten başka bir kullanıcı tarafından alınmış.", http.StatusConflict)
			return
		}
		if emailExists {
			utils.SendJSONError(w, "Bu e-posta adresi zaten başka bir kullanıcı tarafından kullanılıyor.", http.StatusConflict)
			return
		}
	}

	avatarPath := "./static/default-avatar.png"
	file, handler, err := req.FormFile("avatar")
	if err == nil {
		defer file.Close()

		fileExt := strings.ToLower(filepath.Ext(handler.Filename))
		allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
		if !allowedExtensions[fileExt] {
			utils.SendJSONError(w, "Geçersiz fotoğraf formatı. JPG, PNG veya GIF olmalı.", http.StatusBadRequest)
			return
		}

		fileName := uuid.New().String() + fileExt

		uploadDir := "./static/avatars"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err = os.MkdirAll(uploadDir, 0755)
			if err != nil {
				log.Printf("Error creating upload directory %s: %v", uploadDir, err)
				utils.SendJSONError(w, "Sunucu hatası: Avatar yüklenemedi.", http.StatusInternalServerError)
				return
			}
		}

		filePath := filepath.Join(uploadDir, fileName)
		dst, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating file %s: %v", filePath, err)
			utils.SendJSONError(w, "Sunucu hatası: Avatar yüklenemedi.", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Printf("Error copying file to %s: %v", filePath, err)
			utils.SendJSONError(w, "Sunucu hatası: Avatar yüklenemedi.", http.StatusInternalServerError)
			return
		}
		avatarPath = "./static/avatars/" + fileName
	} else if err != http.ErrMissingFile {
		log.Printf("Error getting avatar file: %v", err)
		utils.SendJSONError(w, "Avatar yüklenirken bir hata oluştu.", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec(
		"UPDATE users SET name = $1, username = $2, surname = $3, email = $4, pp = $5 WHERE id = $6",
		name, username, surname, email, avatarPath, userID,
	)

	if err != nil {
		log.Printf("SQL Update Failed for user %s: %v\n", userID, err)
		utils.SendJSONError(w, "Profil güncellenemedi. Lütfen tekrar deneyin.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profil başarıyla güncellendi.",
	})

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
