package entity

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	UserName  string    `json:"userName"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	Pantries  []string  `json:"pantries"` // Pantry IDs
}
