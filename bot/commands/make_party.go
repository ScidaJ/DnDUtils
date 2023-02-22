package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Flow
// Make Role
// Create Channel Category if not exist
// Create Channel if specified
// Give role permission for Channel
// Add Reaction to message, add players as they react
func MakePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, sugar *zap.SugaredLogger) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: makePartyHandler(s, i, sugar),
		},
	})
	message, err := s.InteractionResponse(i.Interaction)

	if err != nil {
		sugar.Panic("error adding reaction ", err)
	}

	err = s.MessageReactionAdd(message.ChannelID, message.ID, "👍")

	if err != nil {
		sugar.Error("error adding reaction ", err)
	}
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

	role, err := CreateNewRole(sugar, s, partyName, i.GuildID)

	if err != nil {
		sugar.Error("error creating role ", err)
		return "Unable to create role"
	}

	if createChannel {
		category, err := CreateNewCategory(sugar, s, i.GuildID)

		if err != nil {
			sugar.Error("error creating category ", err)
			return "Unable to create category"
		}

		_, err = CreateNewChannel(sugar, s, partyName, i.GuildID, role.ID, category.ID)

		if err != nil {
			sugar.Error("error creating channel ", err)
			return "Unable to create party channel"
		}

		return fmt.Sprintf("Created channel and role for party %v", partyName)
	}

	return fmt.Sprintf("Created role for party %v", partyName)
}
