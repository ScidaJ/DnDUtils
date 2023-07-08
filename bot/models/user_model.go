package models

import "strings"

type User struct {
	Id        string   `json:"id,omitempty"`
	DiscordId string   `json:"discord_id,omitempty" validate:"required"`
	Servers   []string `json:"servers,omitempty" validate:"required"`
	Owns      []string `json:"owns,omitempty" validate:"required"`
	Parties   []string `json:"parties,omitempty" validate:"required"`
}

// This is greedy, do not use in production
func PrintAllUsers(s *[]User) string {
	var sb strings.Builder
	for i, e := range *s {
		if i < len(*s)-1 {
			sb.WriteString(e.DiscordId + ", ")
		} else {
			sb.WriteString(e.DiscordId)
		}
	}
	users := sb.String()
	sb.Reset()

	return users
}
