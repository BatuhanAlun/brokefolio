package middleware

import (
	"brokefolio/backend/internal/models" // Assuming your models package is here
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
)

type mailContextKey struct{}

var MailContext = mailContextKey{}

type MailMiddlewareDB struct {
	DB *sql.DB
}

func NewMailMiddlewareDB(db *sql.DB) *MailMiddlewareDB {
	return &MailMiddlewareDB{DB: db}
}

func (h *MailMiddlewareDB) CheckMailParameterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.URL.Query().Get("token")

		if tokenString == "" {
			log.Println("Password reset attempt: No token provided in URL.")
			http.Redirect(w, r, "/unauthorized", http.StatusFound)
			return
		}

		var resetTokenData models.ResetToken

		err := h.DB.QueryRow("SELECT email, expdate FROM resettokens WHERE token = $1", tokenString).Scan(&resetTokenData.Email, &resetTokenData.ExpDate)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("Password reset attempt: Invalid or non-existent token: %s", tokenString)
				http.Redirect(w, r, "/unauthorized", http.StatusFound)
				return
			}
			log.Printf("Database error during token lookup in middleware: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if time.Now().After(resetTokenData.ExpDate) {
			log.Printf("Password reset attempt: Expired token for email %s (token: %s)", resetTokenData.Email, tokenString)

			go func() { //new thread
				_, delErr := h.DB.Exec("DELETE FROM resettokens WHERE token = $1", tokenString)
				if delErr != nil {
					log.Printf("Error deleting expired token %s: %v", tokenString, delErr)
				}
			}()
			http.Redirect(w, r, "/unauthorized", http.StatusFound)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, MailContext, resetTokenData.Email)

		r = r.WithContext(ctx)

		log.Printf("Password reset attempt: Token %s is valid. Proceeding to reset page.", tokenString)
		next.ServeHTTP(w, r)
	})
}
