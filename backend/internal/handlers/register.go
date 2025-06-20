package handlers

import (
	"brokefolio/backend/internal/models"
	"brokefolio/backend/internal/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type RegisterHandler struct {
	DB *sql.DB
}

func NewRegisterHandler(db *sql.DB) *RegisterHandler {
	return &RegisterHandler{DB: db}
}

func (h *RegisterHandler) RegisterHandler(w http.ResponseWriter, req *http.Request) {

	var newUser models.User

	newUser.ID = uuid.New()
	newUser.Role = "user"

	err := json.NewDecoder(req.Body).Decode(&newUser)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println(err)
		return
	}

	defer req.Body.Close()

	if newUser.Email == "" || newUser.Name == "" || newUser.Password == "" || newUser.Surname == "" || newUser.Username == "" || newUser.ConfirmPassword == "" {
		http.Error(w, "Credentials Could not be Empty", http.StatusBadRequest)
		return
	}

	if newUser.Password != newUser.ConfirmPassword {
		http.Error(w, "Passwords do not Match", http.StatusBadRequest)
		return
	}

	hashPass, err := utils.HashPassword(newUser.Password)

	if err != nil {
		utils.SentError(w, "Password Hash Error", http.StatusInternalServerError)
		return
	}

	//email needs to be unique check

	_, err = h.DB.Exec("INSERT INTO users (id, username, password, email, role, name, surname) VALUES ($1,$2,$3,$4,$5,$6,$7)", newUser.ID, newUser.Username, hashPass, newUser.Email, newUser.Role, newUser.Name, newUser.Surname)

	if err != nil {
		http.Error(w, "Could not insert Data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	sucMessage := models.HttpSuccess{
		Message: "Succesfully Served",
	}

	err = json.NewEncoder(w).Encode(sucMessage)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("/register succesfully served POST")
}
