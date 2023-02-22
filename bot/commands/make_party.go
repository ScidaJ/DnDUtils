package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func MakePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, sugar *zap.SugaredLogger) {
	response := ""

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	channelName := ""

	if option, ok := optionMap["party-name"]; ok {
		channelName = option.StringValue()
	}

	createdChannel, err := CreateNewChannel(sugar, channelName, s, i.GuildID)

	if err != nil {
		sugar.Error("error creating channel", err)
	}

	if !createdChannel {
		response = "Unable to create party channel"
	} else {
		response = fmt.Sprintf("Created channel for party %v", channelName)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})

}
