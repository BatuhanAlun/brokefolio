package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type LogOutDBHandler struct {
	DB *sql.DB
}

func NewLogOutHandler(db *sql.DB) *LogOutDBHandler {
	return &LogOutDBHandler{DB: db}
}

func (h *LogOutDBHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request on /logout")

	cookie, err := r.Cookie("sessionID")

	if err != nil {
		http.Redirect(w, r, "/went-wrong", http.StatusFound)
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

	_, _ = h.DB.Exec("DELETE FROM sessions WHERE session_id = $1", cookie.Value)

	log.Println("Session cookie invalidated.")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
