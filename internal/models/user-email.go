package models

import (
	"time"
)

// Email -
type Email struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Email string `json:"email"`
	Verified bool `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserEmail -
func UserEmail() {

}