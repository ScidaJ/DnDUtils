package commands

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var Permissions int64 = 1067403562561
var Mentionable = true

func CreateNewChannel(sugar *zap.SugaredLogger, channelName string, s *discordgo.Session, g string, r string) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(g)

	if err != nil {
		sugar.Error("error fetching channels", err)
		return nil, err
	}

	for _, v := range channels {
		if v.Name == channelName {
			return nil, fmt.Errorf("%v already exists", channelName)
		}
	}

	channel, err := s.GuildChannelCreate(g, channelName, discordgo.ChannelTypeGuildText)

	if err != nil {
		sugar.Error("error creating channel", err)
		return nil, err
	}

	channelPermissions := []*discordgo.PermissionOverwrite{
		{
			ID:    r,
			Type:  0,
			Deny:  0,
			Allow: Permissions,
		},
	}

	channelEditData := discordgo.ChannelEdit{
		PermissionOverwrites: channelPermissions,
	}

	channel, err = s.ChannelEdit(channel.ID, &channelEditData)

	if err != nil {
		sugar.Error("error editing channel", err)
		return nil, err
	}

	return channel, nil
}

func CreateNewRole(sugar *zap.SugaredLogger, roleName string, s *discordgo.Session, g string) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(g)

	if err != nil {
		sugar.Error("error fetching roles", err)
		return nil, fmt.Errorf("unable to create role")
	}

	for _, v := range roles {
		if v.Name == roleName {
			return nil, fmt.Errorf("%v already exists", roleName)
		}
	}

	params := &discordgo.RoleParams{
		Name:        roleName,
		Color:       randomColor(),
		Permissions: &Permissions,
		Mentionable: &Mentionable,
	}

	role, err := s.GuildRoleCreate(g, params)

	if err != nil {
		sugar.Error("error creating role", err)
		return nil, fmt.Errorf("unable to create role")
	}

	return role, nil
}

func randomColor() *int {
	min := 0
	max := 16777215
	random := rand.Intn(max-min+1) + min
	return &random
}
