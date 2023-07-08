package models

// TODO: Image for items?
type Item struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty" validate:"required"`
	Owner string `json:"owner,omitempty" validate:"required"`
}
