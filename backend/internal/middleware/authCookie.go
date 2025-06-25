package middleware

import (
	"brokefolio/backend/internal/utils"
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UsernameContextKey contextKey = "username"
)

func MiddlewareAuthJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("authToken")

		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Authentication cookie Not Found", http.StatusNotFound)
			}
		}

		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, utils.JWTKeyFunc)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		username, ok := claims["username"].(string)

		if !ok {
			http.Error(w, "Invalid token: username not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UsernameContextKey, username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
