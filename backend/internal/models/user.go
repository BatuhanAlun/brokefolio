package models

import "github.com/google/uuid"

type User struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	Name            string    `json:"name"`
	Surname         string    `json:"surname"`
	ConfirmPassword string    `json:"confirm_password"`
	AvatarURL       string    `json:"avatarUrl"`
}
