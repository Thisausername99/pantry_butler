package entity

import (
	"time"
)

type Pantry struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Owner     string         `json:"owner"`
	CreatedAt time.Time      `json:"createdAt"`
	Entries   *[]PantryEntry `json:"entries"`
}

type PantryEntry struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Expiration   *time.Time `json:"expiration,omitempty"`
	Quantity     *int       `json:"quantity,omitempty"`
	QuantityType *string    `json:"quantityType,omitempty"`
}
