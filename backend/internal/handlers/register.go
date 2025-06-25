package handlers

import (
	"brokefolio/backend/internal/utils"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type RegisterHandler struct {
	DB *sql.DB
}

func NewRegisterHandler(db *sql.DB) *RegisterHandler {
	return &RegisterHandler{DB: db}
}

func (h *RegisterHandler) RegisterHandler(w http.ResponseWriter, req *http.Request) {

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
	password := req.FormValue("password")

	if name == "" || surname == "" || username == "" || email == "" || password == "" {
		utils.SendJSONError(w, "Tüm alanlar doldurulmalıdır.", http.StatusBadRequest)
		return
	}
	if len(password) < 8 {
		utils.SendJSONError(w, "Şifre en az 8 karakter olmalıdır.", http.StatusBadRequest)
		return
	}

	var existingUserCount int
	err = h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&existingUserCount)
	if err != nil {
		log.Printf("Error checking existing username: %v", err)
		utils.SendJSONError(w, "Sunucu hatası. Lütfen tekrar deneyin.", http.StatusInternalServerError)
		return
	}
	if existingUserCount > 0 {
		utils.SendJSONError(w, "Bu kullanıcı adı zaten alınmış.", http.StatusConflict)
		return
	}

	err = h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&existingUserCount)
	if err != nil {
		log.Printf("Error checking existing email: %v", err)
		utils.SendJSONError(w, "Sunucu hatası. Lütfen tekrar deneyin.", http.StatusInternalServerError)
		return
	}
	if existingUserCount > 0 {
		utils.SendJSONError(w, "Bu e-posta adresi zaten kayıtlı.", http.StatusConflict)
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		utils.SendJSONError(w, "Şifre işlenirken bir hata oluştu.", http.StatusInternalServerError)
		return
	}

	avatarPath := "../../../static/default-avatar.png"
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

		uploadDir := "../../../static/avatars"
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
		avatarPath = "../../../static/avatars/" + fileName
	} else if err != http.ErrMissingFile {
		log.Printf("Error getting avatar file: %v", err)
		utils.SendJSONError(w, "Avatar yüklenirken bir hata oluştu.", http.StatusBadRequest)
		return
	}

	userID := uuid.New()
	err = h.DB.QueryRow(`
		INSERT INTO users (id , name, surname, username, email, password, pp, role,created_at)
		VALUES ($1, $2, $3, $4, $5, $6,$7, $8, NOW())
		RETURNING id`,
		userID, name, surname, username, email, hashedPassword, avatarPath, "user").Scan(&userID)
	if err != nil {
		log.Printf("Database insert error for new user: %v", err)
		utils.SendJSONError(w, "Kayıt başarısız oldu. Lütfen tekrar deneyin.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kayıt başarılı!"})
}
