package entity

type Recipe struct {
	ID          string                 `json:"id" bson:"_id"`
	Name        string                 `json:"name"`
	Rating      *int                   `json:"rating,omitempty"`
	Ingredients map[string]interface{} `json:"ingredients"`
	Difficulty  *int                   `json:"difficulty,omitempty"`
	Cuisine     *string                `json:"cuisine,omitempty"`
	Description string                 `json:"description"`
	SourceUrl   *string                `json:"source_url,omitempty"`
}
