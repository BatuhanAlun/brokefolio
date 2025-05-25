package main

import (
	"net/http"

	"brokefolio.com/internal/handlers"
)

func main() {

	router := http.NewServeMux()

	router.HandleFunc("POST /api/login", handlers.LoginHandler)
}
