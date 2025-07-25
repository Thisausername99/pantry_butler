package entity

import "time"

type PantryEntry struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Expiration   *time.Time `json:"expiration,omitempty"`
	Quantity     *int       `json:"quantity,omitempty"`
	QuantityType *string    `json:"quantityType,omitempty"`
}
