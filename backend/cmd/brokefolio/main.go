package main

import (
	"brokefolio/backend/internal/database"
	"brokefolio/backend/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {

	db := database.InitDB()
	defer db.Close()

	router := http.NewServeMux()

	registerHandler := handlers.NewRegisterHandler(db)

	router.HandleFunc("POST /api/login", handlers.LoginHandler)
	router.HandleFunc("/api/register", registerHandler.RegisterHandler)

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://127.0.0.1:5500"
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{allowedOrigins},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          os.Getenv("CORS_DEBUG") == "true",
	})

	handler := c.Handler(router)

	log.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
