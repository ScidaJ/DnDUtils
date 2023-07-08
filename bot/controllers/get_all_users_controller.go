package controllers

import (
	"dndutils/bot/configs"
	"dndutils/bot/models"
	"dndutils/bot/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func GetAllUsersHandler(s *discordgo.Session, i *discordgo.InteractionCreate, e string, sugar *zap.SugaredLogger) {
	response, err := getAllUsersHandler(i, e, sugar)
	if err != nil {
		sugar.Error("error getting all users", err)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func getAllUsersHandler(i *discordgo.InteractionCreate, e string, sugar *zap.SugaredLogger) (string, error) {
	var users models.GetAllUsersResponse
	approvedList := configs.EnvApprovedUsers()

	if !utils.IncludesString(approvedList, i.Member.User.ID) {
		return "", errors.New("error getting all users")
	}

	uri, err := url.JoinPath(e, UserRoute)
	if err != nil {
		sugar.Error(err)
	}

	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		sugar.Error("error reading api response ", err)
		return "", err
	}

	err = json.Unmarshal(bodyBytes, &users)
	if err != nil {
		sugar.Error("error unmarshalling api response ", err)
		return "", err
	}

	return models.PrintAllUsers(&users.Data.Data), nil
}
