package models

import (
	"time"
)

// Profile -
type Profile struct {
	UserID int `json:"user_id"`
	ProfileImage string `json:"profile_image"`
	Bio string `json:"bio"`
	TeamID int `json:"team_id"`
	Settings string `json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Profiles -
func Profiles () {

}
