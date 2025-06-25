package handlers

import (
	"brokefolio/backend/internal/models"
	"brokefolio/backend/internal/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type LoginDBHandler struct {
	DB *sql.DB
}

func NewLoginHandler(db *sql.DB) *LoginDBHandler {
	return &LoginDBHandler{DB: db}
}

func (h *LoginDBHandler) LoginHandler(w http.ResponseWriter, req *http.Request) {

	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	if user.Username == "" && user.Password == "" {
		http.Error(w, "Credentials shouldn't be empty!", http.StatusBadRequest)
		return
	}

	var sqlUsername string
	var sqlPassword string
	var sqlRole string
	var sqlUserID string

	err = h.DB.QueryRow("SELECT id ,username , password , role FROM users WHERE username = $1", user.Username).Scan(&sqlUserID, &sqlUsername, &sqlPassword, &sqlRole)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	token, err := utils.GenerateJWT(user.Username, sqlRole)

	if err != nil {
		http.Error(w, "Error generating JWT", http.StatusInternalServerError)
		return
	}

	err = utils.CheckPasswordHash(user.Password, sqlPassword)
	if err != nil {
		http.Error(w, "Wrong Password!", http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     "authToken",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	//session

	sessionID, err := utils.GenerateNewToken(32)

	if err != nil {
		http.Error(w, "Error generating session ID", http.StatusInternalServerError)
		log.Printf("Failed to generate session id %v\n", err)
		return
	}

	sessionExpiresAt := time.Now().Add(24 * time.Hour)

	_, err = h.DB.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES ($1, $2, $3)", sessionID, sqlUserID, sessionExpiresAt)

	if err != nil {
		http.Error(w, "Error Saving Session", http.StatusInternalServerError)
		log.Printf("Failed to insert session into DB %v\n", err)
		return
	}

	sessionCookie := &http.Cookie{
		Name:     "sessionID",
		Value:    sessionID,
		Expires:  sessionExpiresAt,
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, sessionCookie)

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Succesfully Logged In!",
	})
	log.Println("Login succesful, JWT sent")
	log.Println("Session Token Sent")

}
