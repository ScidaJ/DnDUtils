package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var Permissions int64 = 1067403562561
var BlockChannel int64 = 1024
var Mentionable = true
var CategoryName = "Campaign Channels"
var ReactionTimeout = time.Duration(300 * float64(time.Second))
var Emoji = "👍"
var PostRoute = "party"

// TODO: Send party/player info to API for DB party/player creation
// Handler for /make-party command. Adds reacitons
// for users to join party
func MakePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, e string, sugar *zap.SugaredLogger) {
	response, role, err := makePartyHandler(s, i, e, sugar)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})

	if err != nil {
		sugar.Error("error making party ", err)
		return
	}

	message, err := s.InteractionResponse(i.Interaction)

	if err != nil {
		sugar.Error("error getting parent message ", err)
		return
	}

	err = s.MessageReactionAdd(message.ChannelID, message.ID, Emoji)

	if err != nil {
		sugar.Error("error adding reaction ", err)
		return
	}

	// Below adapted from https://github.com/Necroforger/dgwidgets
	startTime := time.Now()

	var reaction *discordgo.MessageReaction

	for {
		select {
		case k := <-nextMessageReactionAddC(s):
			reaction = k.MessageReaction
		case <-time.After(time.Until(startTime.Add(ReactionTimeout))):
			s.MessageReactionsRemoveAll(message.ChannelID, message.ID)
			return
		}

		if reaction.MessageID != message.ID || s.State.User.ID == reaction.UserID {
			continue
		}

		err := s.GuildMemberRoleAdd(reaction.GuildID, reaction.UserID, role)

		if err != nil {
			sugar.Error("error adding role to user ", err)
			return
		}

		go func() {
			time.Sleep(time.Millisecond * 250)
			s.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, Emoji, reaction.UserID)
		}()
	}
}

// Handler for the make-party command
func makePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, e string, sugar *zap.SugaredLogger) (string, string, error) {

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

	err := partyAPICall(partyName, i.GuildID, e, sugar)

	role, err := CreateNewRole(sugar, s, partyName, i.GuildID)

	if err != nil {
		sugar.Error("error creating role ", err)
		return "Unable to create role", "", err
	}

	if createChannel {
		category, err := CreateNewCategory(sugar, s, i.GuildID)

		if err != nil {
			sugar.Error("error creating category ", err)
			return "Unable to create category", "", err
		}

		_, err = CreateNewChannel(sugar, s, partyName, i.GuildID, role.ID, category.ID)

		if err != nil {
			sugar.Error("error creating channel ", err)
			return "Unable to create party channel", "", err
		}

		return fmt.Sprintf("Created channel and role for party %v. React to this message to join it.", partyName), role.ID, nil
	}

	return fmt.Sprintf("Created role for party %v. React to this message to join it.", partyName), role.ID, nil
}

// Creates new text channel with a given name
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

func partyAPICall(partyName string, serverID string, endpoint string, sugar *zap.SugaredLogger) error {
	uri, err := url.JoinPath(endpoint, PostRoute)

	if err != nil {
		sugar.Error("error building url ", err)
		return err
	}

	postBody, _ := json.Marshal(map[string]string{
		"party_name": partyName,
		"server_id":  serverID,
	})

	responseBody := bytes.NewBuffer(postBody)

	_, err = http.Post(uri, "application/json", responseBody)

	if err != nil {
		sugar.Error("error creating party with api ", err)
		return err
	}

	return nil
}
