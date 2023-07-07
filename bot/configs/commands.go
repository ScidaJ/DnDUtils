package configs

import (
	"dndutils/bot/models"

	"github.com/bwmarrin/discordgo"
)

// Different commands that the bot responds to. Visit the documentation to see them all.
var (
	MakeParty = models.SlashCommand{
		Name:        "make-party",
		Description: "Make a channel and gives players permissions to view the channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "party-name",
				Description: "Name of the new party. Must be unique to the server",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "make-channel",
				Description: "Make a new channel for the party?",
				Required:    true,
			},
		},
	}
)
