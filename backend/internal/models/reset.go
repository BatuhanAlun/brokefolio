package models

import (
	"time"
)

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type ResetToken struct {
	Token   string    `json:"token"`
	Email   string    `json:"email"`
	ExpDate time.Time `json:"expdate"`
}

type ChangePasswordMailReq struct {
	NewPassword string `json:"newPassword"`
	UserEmail   string `json:"email"`
}
