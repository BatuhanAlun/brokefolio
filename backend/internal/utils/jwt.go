package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(username, role string) (string, error) {
	var jwtEnv = os.Getenv("JWT_SECRET_KEY")

	if jwtEnv == "" {
		log.Fatalf("FATAL ERROR JWT KEY NOT FOUND")
	}

	var jwtKey = []byte(jwtEnv)

	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 720).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTKeyFunc(token *jwt.Token) (interface{}, error) {

	var jwtEnv = os.Getenv("JWT_SECRET_KEY")

	if jwtEnv == "" {
		log.Fatalf("FATAL ERROR JWT KEY NOT FOUND")
	}

	var jwtKey = []byte(jwtEnv)
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return jwtKey, nil
}
