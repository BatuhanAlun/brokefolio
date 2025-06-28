package handlers

import (
	"brokefolio/backend/internal/models"
	"brokefolio/backend/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type MailHandler struct {
	DB *sql.DB
}

func NewMailHandler(db *sql.DB) *MailHandler {

	return &MailHandler{DB: db}

}

func (h *MailHandler) MailResetHandler(w http.ResponseWriter, req *http.Request) {

	var user models.User

	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {

		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	_, err = h.DB.Exec("SELECT email from users WHERE email = $1", user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		} else {
			http.Error(w, "SQL Error", http.StatusInternalServerError)
			return
		}
	}

	resetToken, err := utils.GenerateNewToken(32)

	if err != nil {
		http.Error(w, "Token Generation Failed", http.StatusInternalServerError)
	}

	_, err = h.DB.Exec("INSERT INTO resettokens (token, email, expdate) VALUES ($1,$2,$3)", resetToken, user.Email, time.Now().Add(1*time.Hour))

	if err != nil {
		http.Error(w, "SQL error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	resetLink := fmt.Sprintf("http://localhost:8080/reset-password?token=%s", resetToken)
	subject := "Şifre Resetleme"
	body := fmt.Sprintf("Şifre yenileme isteğiniz:\n\n%s\n\nYukarıdaki linkin geçerlilik süresi 1 saattir.", resetLink)

	err = utils.SendEmail(user.Email, subject, body)
	if err != nil {
		http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
		log.Printf("Error sending email: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Gelen Kutunuzu Kontrol Ediniz!",
	})

	log.Println("Succesfully Served /api/passwordRecovery")

}

func (h *MailHandler) PasswordResetUsingMailHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("Received request on /api/user/change-password-mail")
	var cpReq models.ChangePasswordMailReq

	err := json.NewDecoder(req.Body).Decode(&cpReq)

	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		log.Printf("Invalid Request %v\n", err)
		return
	}

	fmt.Println(cpReq.UserEmail)

	if len(cpReq.NewPassword) < 8 {
		http.Error(w, "Yeni şifre en az 8 karakter olmalıdır.", http.StatusBadRequest)
		return
	}

	var hashedPasswordFromDB string
	err = h.DB.QueryRow("SELECT password FROM users WHERE email = $1", cpReq.UserEmail).Scan(&hashedPasswordFromDB)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Kullanıcı bulunamadı (auth hatası).", http.StatusUnauthorized)
		} else {
			log.Printf("DB error fetching password for user %s: %v", cpReq.UserEmail, err)
			http.Error(w, "Sunucu hatası: Şifre doğrulanamadı.", http.StatusInternalServerError)
		}
		return
	}

	newHashedPassword, err := utils.HashPassword(cpReq.NewPassword)

	if err != nil {
		log.Printf("Bcrypt hashing error for user %s: %v", cpReq.UserEmail, err)
		http.Error(w, "Sunucu hatası: Yeni şifre işlenemedi.", http.StatusInternalServerError)
		return
	}

	_, err = h.DB.Exec("UPDATE users SET password = $1 WHERE email = $2", newHashedPassword, cpReq.UserEmail)
	if err != nil {
		log.Printf("DB error updating password for user %s: %v", cpReq.UserEmail, err)
		http.Error(w, "Sunucu hatası: Şifre güncellenemedi.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Şifreniz başarıyla değiştirildi."})

}
