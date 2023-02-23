package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var Permissions int64 = 1067403562561
var BlockChannel int64 = 1024
var Mentionable = true
var CategoryName = "Campaign Channels"

// TODO: Give role to users upon reaction
// Flow
// Make Role
// Create Channel Category if not exist
// Create Channel if specified
// Give role permission for Channel
// Add Reaction to message, add players as they react
func MakePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, sugar *zap.SugaredLogger) {
	response, err := makePartyHandler(s, i, sugar)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})

	if err != nil {
		sugar.Panic("error making party ", err)
		return
	}

	message, err := s.InteractionResponse(i.Interaction)

	if err != nil {
		sugar.Panic("error adding reaction ", err)
		return
	}

	err = s.MessageReactionAdd(message.ChannelID, message.ID, "👍")

	if err != nil {
		sugar.Error("error adding reaction ", err)
		return
	}
}

func makePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, sugar *zap.SugaredLogger) (string, error) {
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
		return "Unable to create role", err
	}

	if createChannel {
		category, err := CreateNewCategory(sugar, s, i.GuildID)

		if err != nil {
			sugar.Error("error creating category ", err)
			return "Unable to create category", err
		}

		_, err = CreateNewChannel(sugar, s, partyName, i.GuildID, role.ID, category.ID)

		if err != nil {
			sugar.Error("error creating channel ", err)
			return "Unable to create party channel", err
		}

		return fmt.Sprintf("Created channel and role for party %v", partyName), nil
	}

	return fmt.Sprintf("Created role for party %v", partyName), nil
}

func CreateNewChannel(sugar *zap.SugaredLogger, s *discordgo.Session, channelName string, g string, r string, c string) (*discordgo.Channel, error) {
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

	err = s.ChannelPermissionSet(channel.ID, r, discordgo.PermissionOverwriteTypeRole, Permissions, 0)

	if err != nil {
		sugar.Error("error adding role permissions to channel ", err)
		return nil, err
	}

	err = s.ChannelPermissionSet(channel.ID, s.State.User.ID, discordgo.PermissionOverwriteTypeMember, Permissions, 0)

	if err != nil {
		sugar.Error("error adding bot permissions to channel ", err)
		return nil, err
	}

	err = s.ChannelPermissionSet(channel.ID, g, discordgo.PermissionOverwriteTypeRole, 0, BlockChannel)

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

// Creates a new role with the given role name.
// Flow
// Fetches all roles
// Searches through them to see if role exists. If it does, return error
// If it does not, then with the set role params, create the role.
// If an error is returned, return the error
// If no error is returned, return the role
func CreateNewRole(sugar *zap.SugaredLogger, s *discordgo.Session, roleName string, g string) (*discordgo.Role, error) {
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
		Color:       randomColor(),
		Permissions: &Permissions,
		Mentionable: &Mentionable,
	}

	role, err := s.GuildRoleCreate(g, params)

	if err != nil {
		sugar.Error("error creating role ", err)
		return nil, fmt.Errorf("unable to create role")
	}

	return role, nil
}

func CreateNewCategory(sugar *zap.SugaredLogger, s *discordgo.Session, g string) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(g)

	if err != nil {
		sugar.Error("error fetching channels ", err)
		return nil, err
	}

	for _, v := range channels {
		if v.Name == CategoryName {
			return v, nil
		}
	}

	category, err := s.GuildChannelCreate(g, CategoryName, discordgo.ChannelTypeGuildCategory)

	if err != nil {
		sugar.Error("error creating channel ", err)
		return nil, err
	}

	return category, nil
}
