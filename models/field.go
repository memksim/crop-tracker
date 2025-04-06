package models

type Field struct {
	ID     int     `json:"id" db:"id"`
	Name   string  `json:"name" db:"name"`
	AreaHa float64 `json:"area_ha" db:"area_ha"`
	Region string  `json:"region" db:"region"`
}
