package commands

import "github.com/bwmarrin/discordgo"

type SlashCommand struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
}

var (
	MakeParty = SlashCommand{
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
