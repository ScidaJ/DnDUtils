package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Party struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	ServerId string             `json:"server_id,omitempty" validate:"required"`
	Owner    string             `json:"owner,omitempty" validate:"required"`
	Users    []string           `json:"users,omitempty" validate:"required"`
}
