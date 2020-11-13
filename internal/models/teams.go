package models

import (
	"time"
)

// Team -
type Team struct {
	ID int `json:"id"`
	Name int `json:"name"`
	Description string `json:"description"`
	TeamInvite string `json:"team_invite"`
	TeamLogo string `json:"team_logo"`
	TeamSettings string `json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Teams -
func Teams() {

}