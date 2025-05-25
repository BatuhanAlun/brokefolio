package handlers

import (
	"encoding/json"
	"net/http"

	"brokefolio.com/internal/models"
)

func LoginHandler(w http.ResponseWriter, req *http.Request) {

	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

}
