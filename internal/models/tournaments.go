package models

import (
	"time"
)

// Tournament -
type Tournament struct {
	ID int `json:"id"`
	TeamID int `json:"team_id"`
	SeasonType string `json:"season_type"`
	TournamentName string `json:"tournament_name"`
	Description string `json:"description"`
	Settings string `json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Tournaments -
func Tournaments() {

}