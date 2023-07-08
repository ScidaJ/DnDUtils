package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"dndutils/bot/models"
	"dndutils/bot/utils"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var Permissions int64 = 1067403562561
var HideChannel int64 = 1024
var Mentionable = true
var CategoryName = "Campaign Channels"
var ReactionTimeout = time.Duration(300 * float64(time.Second))
var Emoji = "👍"
var PostPartyRoute = "party"

// Handler for /make-party command. Adds reacitons
// for users to join party
func MakePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, e string, sugar *zap.SugaredLogger) {
	response, role, partyId, err := makePartyHandler(s, i, e, sugar)

	if err != nil {
		sugar.Error("error making party ", err)
		return
	}

	// Adding the original author as a user
	err = AddUserHandler(i.Member.User.ID, i.GuildID, partyId, e)

	if err != nil {
		sugar.Error("error adding user ", err)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})

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

	go func() {
		for {
			select {
			case k := <-utils.NextMessageReactionAddC(s):
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

			err = AddUserHandler(reaction.UserID, reaction.GuildID, partyId, e)

			if err != nil {
				sugar.Error("error adding user ", err)
			}

			go func() {
				time.Sleep(time.Millisecond * 250)
				err := s.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, Emoji, reaction.UserID)
				if err != nil {
					sugar.Error("error removing reaction", err)
				}
			}()
		}
	}()
}

// Handler for the make-party command
func makePartyHandler(s *discordgo.Session, i *discordgo.InteractionCreate, e string, sugar *zap.SugaredLogger) (string, string, string, error) {

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

	partyId, err := partyAPICall(partyName, i.GuildID, i.Member.User.ID, e, sugar)

	if err != nil {
		sugar.Error("error sending party to api ", err)
		return "Unable to create party", "", "", err
	}

	role, err := utils.CreateNewRole(s, i.GuildID, Permissions, Mentionable, partyName, sugar)

	if err != nil {
		sugar.Error("error creating role ", err)
		return "Unable to create role", "", "", err
	}

	if createChannel {
		category, err := utils.CreateNewCategory(s, i.GuildID, CategoryName, sugar)

		if err != nil {
			sugar.Error("error creating category ", err)
			return "Unable to create category", "", "", err
		}

		_, err = utils.CreateNewChannel(s, i.GuildID, role.ID, category.ID, Permissions, HideChannel, partyName, sugar)

		if err != nil {
			sugar.Error("error creating channel ", err)
			return "Unable to create party channel", "", "", err
		}

		return fmt.Sprintf("Created channel and role for party %v. React to this message to join it.", partyName), role.ID, "", nil
	}

	return fmt.Sprintf("Created role for party %v. React to this message to join it.", partyName), role.ID, partyId, nil
}

func partyAPICall(partyName string, serverID string, owner string, endpoint string, sugar *zap.SugaredLogger) (string, error) {
	uri, err := url.JoinPath(endpoint, PostPartyRoute)

	if err != nil {
		sugar.Error("error building url ", err)
		return "", err
	}

	users := []string{
		owner,
	}

	postBody, _ := json.Marshal(models.Party{
		Name:     partyName,
		ServerId: serverID,
		Owner:    owner,
		Users:    users,
	})

	payload := bytes.NewBuffer(postBody)

	resp, err := http.Post(uri, "application/json", payload)
	if err != nil {
		sugar.Error("error creating party with api ", err)
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		sugar.Error("error reading api response ", err)
		return "", err
	}

	var data models.PostPartyResponse

	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		sugar.Error("error unmarshalling api response ", err)
		return "", err
	}

	return data.Data.Data.InsertedID, nil
}
