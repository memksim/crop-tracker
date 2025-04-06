package models

type Sowing struct {
	ID      int    `json:"id" db:"id"`
	FieldID int    `json:"field_id" db:"field_id"`
	Crop    string `json:"crop" db:"crop"`
	SowedAt string `json:"sowed_at" db:"sowed_at"` // Можно использовать time.Time, но строка проще
}
