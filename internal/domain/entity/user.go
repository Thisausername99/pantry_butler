package entity

import (
	"time"
)

type User struct {
	ID        string    `json:"id" bson:"id"`
	UserName  string    `json:"userName" bson:"userName"`
	Password  string    `json:"password" bson:"password"`
	Email     string    `json:"email" bson:"email"`
	FirstName string    `json:"firstName" bson:"firstName"`
	LastName  string    `json:"lastName" bson:"lastName"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Pantries  []string  `json:"pantries" bson:"pantries"` // Pantry IDs
}
