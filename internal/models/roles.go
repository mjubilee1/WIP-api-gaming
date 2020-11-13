package models

import (
	"time"
)

// Role -
type Role struct {
	ID int `json:"id"`
	Role int `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Roles -
func Roles () {

}
