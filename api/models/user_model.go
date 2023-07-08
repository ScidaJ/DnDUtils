package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_idomitempty" validate:"required"`
	DiscordId string             `json:"discord_id,omitempty" validate:"required"`
	Servers   []string           `json:"servers,omitempty" validate:"required"`
	Owns      []string           `json:"owns,omitempty" validate:"required"`
	Parties   []string           `json:"parties,omitempty" validate:"required"`
}
