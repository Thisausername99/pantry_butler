package entity

type User struct {
	ID        string        `json:"id"`
	Email     string        `json:"email"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Pantry    []PantryEntry `json:"pantry"`
}
