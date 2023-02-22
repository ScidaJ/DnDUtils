package commands

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func CreateNewChannel(sugar *zap.SugaredLogger, channelName string, s *discordgo.Session, g string) (bool, error) {
	channels, err := s.GuildChannels(g)

	if err != nil {
		sugar.Error("error creating channel", err)
		return false, err
	}

	for _, v := range channels {
		if v.Name == channelName {
			return false, nil
		}
	}

	_, err = s.GuildChannelCreate(g, channelName, discordgo.ChannelTypeGuildText)

	if err != nil {
		sugar.Error("error creating channel", err)
		return false, err
	}

	return true, nil
}
