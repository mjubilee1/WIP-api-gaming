package models;

import (
	"time"
)

// League -
type League struct {
	ID int64 `json:"id"`
	TeamID int `json:"team_id"`
	SeasonTypeID string `json:"season_type_id"`
	LeagueName string `json:"league_name"`
	Description string `json:"description"`
	LeagueSettings string `json:"league_settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Leagues -
func Leagues() {

}