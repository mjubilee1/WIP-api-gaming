package models

import (
	"time"
)

// GameMode -
type GameMode struct {
	ID int64 `json:"id"`
	GameModeName int `json:"game_mode_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GameModes -
func GameModes () {

}
