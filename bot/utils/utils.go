package utils

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Gets a random color for a party role
func RandomColor() *int {
	min := 0
	max := 16777215
	random := rand.Intn(max-min+1) + min
	return &random
}

// Taken from https://github.com/Necroforger/dgwidgets/blob/master/util.go#L16-L23
func NextMessageReactionAddC(s *discordgo.Session) chan *discordgo.MessageReactionAdd {
	out := make(chan *discordgo.MessageReactionAdd)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageReactionAdd) {
		out <- e
	})
	return out
}

// Loops over a slice s to check if it contains key k.
func IncludesString(s []string, k string) bool {
	for _, v := range s {
		if v == k {
			return true
		}
	}

	return false
}

/*
Creates new text channel with a given name
Params
s - Discord Session
g - Guild/Server ID
r - Role ID to apply to channel
c - Category ID for channel to go under
p - Permission int64 for channel
b - Permission int64 which does not allow non-members to view channel
channelName - Name of the channel
sugar - TODO Replace with Zerolog

Returns
Pointer to channel object
Error
*/
func CreateNewChannel(s *discordgo.Session, g string, r string, c string, p int64, b int64, channelName string, sugar *zap.SugaredLogger) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(g)

	if err != nil {
		sugar.Error("error fetching channels ", err)
		return nil, err
	}

	for _, v := range channels {
		if v.Name == channelName {
			return nil, fmt.Errorf("%v already exists", channelName)
		}
	}

	channel, err := s.GuildChannelCreate(g, channelName, discordgo.ChannelTypeGuildText)

	if err != nil {
		sugar.Error("error creating channel ", err)
		return nil, err
	}

	err = s.ChannelPermissionSet(channel.ID, r, discordgo.PermissionOverwriteTypeRole, p, 0)

	if err != nil {
		sugar.Error("error adding role permissions to channel ", err)
		return nil, err
	}

	err = s.ChannelPermissionSet(channel.ID, s.State.User.ID, discordgo.PermissionOverwriteTypeMember, p, 0)

	if err != nil {
		sugar.Error("error adding bot permissions to channel ", err)
		return nil, err
	}

	err = s.ChannelPermissionSet(channel.ID, g, discordgo.PermissionOverwriteTypeRole, 0, b)

	if err != nil {
		sugar.Error("error adding everyone permissions to channel ", err)
		return nil, err
	}

	channelEditData := discordgo.ChannelEdit{
		ParentID: c,
	}

	channel, err = s.ChannelEdit(channel.ID, &channelEditData)

	if err != nil {
		sugar.Error("error editing channel ", err)
		return nil, err
	}

	return channel, nil
}

/*
Creates a new role with given params
Params
s - Discord Session
g - Guild/Server ID
p - Permission int64 for channel
m - Boolean value to set if role is mentionable
roleName - Name of the channel
sugar - TODO Replace with Zerolog

Returns
Pointer to channel object
Error
*/
func CreateNewRole(s *discordgo.Session, g string, p int64, m bool, roleName string, sugar *zap.SugaredLogger) (*discordgo.Role, error) {
	roles, err := s.GuildRoles(g)

	if err != nil {
		sugar.Error("error fetching roles ", err)
		return nil, fmt.Errorf("unable to create role")
	}

	for _, v := range roles {
		if v.Name == roleName {
			return nil, fmt.Errorf("%v already exists", roleName)
		}
	}

	params := &discordgo.RoleParams{
		Name:        roleName,
		Color:       RandomColor(),
		Permissions: &p,
		Mentionable: &m,
	}

	role, err := s.GuildRoleCreate(g, params)

	if err != nil {
		sugar.Error("error creating role ", err)
		return nil, fmt.Errorf("unable to create role")
	}

	return role, nil
}

/*
Creates a new category for channels with the given params. Optional
Params
s - Discord Session
g - Guild/Server ID
categoryName - Name of the channel
sugar - TODO Replace with Zerolog

Returns
Pointer to channel object
Error
*/
func CreateNewCategory(s *discordgo.Session, g string, categoryName string, sugar *zap.SugaredLogger) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(g)

	if err != nil {
		sugar.Error("error fetching channels ", err)
		return nil, err
	}

	for _, v := range channels {
		if v.Name == categoryName {
			return v, nil
		}
	}

	category, err := s.GuildChannelCreate(g, categoryName, discordgo.ChannelTypeGuildCategory)

	if err != nil {
		sugar.Error("error creating channel ", err)
		return nil, err
	}

	return category, nil
}
