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

	resetLink := fmt.Sprintf("http://yourfrontenddomain.com/reset-password?token=%s", resetToken)
	subject := "Password Reset Request"
	body := fmt.Sprintf("You have requested a password reset. Please click the following link to reset your password:\n\n%s\n\nThis link will expire in 1 hour.", resetLink)

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
