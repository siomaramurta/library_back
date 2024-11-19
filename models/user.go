package models

import (
	"time"
)

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    *int64 `json:"updated_at"`
	IsDeleted    bool   `json:"is_deleted"`
	DeletedAt    *int64 `json:"deleted_at"`
}

func NewUser(name, email, password_hash string) User {
	return User{
		Name:         name,
		Email:        email,
		PasswordHash: password_hash,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    nil,
		IsDeleted:    false,
		DeletedAt:    nil,
	}
}

func (u *User) UpdateTimestamps(isDeleted bool) {
	now := time.Now().Unix()
	u.UpdatedAt = &now
	if isDeleted {
		u.IsDeleted = true
		u.DeletedAt = &now
	}
}
