package models

type Party struct {
	Id       string   `json:"id,omitempty"`
	Name     string   `json:"name,omitempty" validate:"required"`
	ServerId string   `json:"server_id,omitempty" validate:"required"`
	Owner    string   `json:"owner,omitempty" validate:"required"`
	Users    []string `json:"users,omitempty" validate:"required"`
}
