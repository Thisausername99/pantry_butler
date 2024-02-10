package model

type PantryEntry struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Expiration *int   `json:"expiration,omitempty"`
	Quantity   int    `json:"quantity"`
}
