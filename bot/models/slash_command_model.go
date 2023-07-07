package models

import "github.com/bwmarrin/discordgo"

type SlashCommand struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
}
