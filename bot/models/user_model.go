package models

type User struct {
	Id      string   `json:"id,omitempty" validate:"required"`
	Servers []string `json:"servers,omitempty" validate:"required"`
	Owns    []string `json:"owns,omitempty" validate:"required"`
	Parties []string `json:"parties,omitempty" validate:"required"`
}
