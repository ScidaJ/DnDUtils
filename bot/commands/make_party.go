package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Flow
// Make Role
// Create Channel is specified
// Give role permission for Channel
// Add Reaction to message, add players as they react
func MakePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, sugar *zap.SugaredLogger) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: makePartyHandler(s, i, sugar),
		},
	})

}

func makePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, sugar *zap.SugaredLogger) string {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	partyName := ""

	if option, ok := optionMap["party-name"]; ok {
		partyName = option.StringValue()
	}

	createChannel := false

	if option, ok := optionMap["make-channel"]; ok {
		createChannel = option.BoolValue()
	}

	role, err := CreateNewRole(sugar, partyName, s, i.GuildID)

	if err != nil {
		sugar.Error("error creating role", err)
		return "Unable to create role"
	}

	if createChannel {

		_, err := CreateNewChannel(sugar, partyName, s, i.GuildID, role.ID)

		if err != nil {
			sugar.Error("error creating channel", err)
			return "Unable to create party channel"
		}

		return fmt.Sprintf("Created channel and role for party %v", partyName)
	}

	return fmt.Sprintf("Created role for party %v", partyName)
}
