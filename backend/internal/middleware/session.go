package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"
)

type sessionContextKey struct{}
type userContextKey struct{}
type isAuthenticatedContextKey struct{}

var (
	SessionIDContextKey       = sessionContextKey{}
	UserIDContextKey          = userContextKey{}
	IsAuthenticatedContextKey = isAuthenticatedContextKey{}
)

type SessionMiddlewareDB struct {
	DB *sql.DB
}

func NewSessionMiddleware(db *sql.DB) *SessionMiddlewareDB {
	return &SessionMiddlewareDB{DB: db}
}

func (h *SessionMiddlewareDB) CheckSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		isAuthenticated := false
		var sessionID string
		var userID string

		cookie, err := r.Cookie("sessionID")
		if err == nil {
			sessionID = cookie.Value

			var dbUserID string
			var dbExpiresAt time.Time

			err = h.DB.QueryRow(
				"SELECT user_id, expires_at FROM sessions WHERE session_id = $1",
				sessionID,
			).Scan(&dbUserID, &dbExpiresAt)

			if err != nil {
				if err == sql.ErrNoRows {

					log.Printf("Session ID '%s' not found in database.", sessionID)
					ClearSessionCookie(w)
				} else {

					log.Printf("Database lookup error for session ID '%s': %v", sessionID, err)

				}

			} else {

				if time.Now().Before(dbExpiresAt) {

					isAuthenticated = true
					userID = dbUserID
					log.Printf("Session ID '%s' found and valid for UserID: %s", sessionID, userID)
				} else {
					log.Printf("Session ID '%s' found but expired for UserID: %s. Deleting from DB.", sessionID, dbUserID)
					_, _ = h.DB.Exec("DELETE FROM sessions WHERE session_id = $1", sessionID)
					ClearSessionCookie(w)
				}
			}

		} else if err == http.ErrNoCookie {

			log.Println("SessionID cookie not found in request.")
		} else {

			log.Printf("Error retrieving sessionID cookie: %v", err)
		}

		ctx = context.WithValue(ctx, IsAuthenticatedContextKey, isAuthenticated)
		ctx = context.WithValue(ctx, UserIDContextKey, userID)
		ctx = context.WithValue(ctx, SessionIDContextKey, sessionID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}
