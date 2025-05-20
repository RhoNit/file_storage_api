package models

import (
	"time"
)

type User struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	StorageUsed  int64     `json:"storageUsed"`
	CreatedAt    time.Time `json:"createdAt"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
