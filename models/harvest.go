package models

type Harvest struct {
	ID         int     `json:"id" db:"id"`
	FieldID    int     `json:"field_id" db:"field_id"`
	Crop       string  `json:"crop" db:"crop"`
	YieldPerHa float64 `json:"yield_t_per_ha" db:"yield_t_per_ha"`
}
