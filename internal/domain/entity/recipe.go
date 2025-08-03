package entity

type Recipe struct {
	ID          string                 `json:"id" bson:"id"`
	Name        string                 `json:"name" bson:"name"`
	Rating      *int                   `json:"rating,omitempty" bson:"rating,omitempty"`
	Ingredients map[string]interface{} `json:"ingredients" bson:"ingredients"`
	Difficulty  *int                   `json:"difficulty,omitempty" bson:"difficulty,omitempty"`
	Cuisine     *string                `json:"cuisine,omitempty" bson:"cuisine,omitempty"`
	Description string                 `json:"description" bson:"description"`
	SourceUrl   *string                `json:"source_url,omitempty" bson:"source_url,omitempty"`
}
