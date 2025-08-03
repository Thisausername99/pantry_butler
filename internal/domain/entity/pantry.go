package entity

import (
	"time"
)

type Pantry struct {
	ID        string         `json:"id" bson:"id"`
	Name      string         `json:"name" bson:"name"`
	OwnerID   string         `json:"ownerId" bson:"ownerId"`
	CreatedAt time.Time      `json:"createdAt" bson:"createdAt"`
	Entries   *[]PantryEntry `json:"pantry_entries" bson:"pantry_entries"`
}

type PantryEntry struct {
	ID           string     `json:"id" bson:"id"`
	Name         string     `json:"name" bson:"name"`
	Expiration   *time.Time `json:"expiration,omitempty" bson:"expiration,omitempty"`
	Quantity     *float64   `json:"quantity,omitempty" bson:"quantity,omitempty"`
	QuantityType *string    `json:"quantityType,omitempty" bson:"quantityType,omitempty"`
}
