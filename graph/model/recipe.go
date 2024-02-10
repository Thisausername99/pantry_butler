package model

type Recipe struct {
	Rating      int      `json:"rating"`
	Ingredients []string `json:"ingredients"`
	Difficulty  int      `json:"difficulty"`
	Cuisine     string   `json:"cuisine"`
	Description string   `json:"description"`
}
