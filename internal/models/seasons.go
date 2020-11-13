package models

import (
	"time"
)

// Season -
type Season struct {
	ID int `json:"id"`
	SeasonID int `json:"season_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Seasons -
func Seasons() {

}