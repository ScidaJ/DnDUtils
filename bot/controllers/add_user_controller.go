package controllers

import (
	"bytes"
	"dndutils/bot/models"
	"encoding/json"
	"net/http"
	"net/url"
)

var UserRoute = "user"

func AddUserHandler(userId string, guildId string, partyId string, endpoint string) error {
	uri, err := url.JoinPath(endpoint, UserRoute)

	if err != nil {
		//sugar.Error("error building url ", err)
		return err
	}

	servers := []string{
		guildId,
	}

	parties := []string{
		partyId,
	}

	postBody, _ := json.Marshal(models.User{
		DiscordId: userId,
		Servers:   servers,
		Parties:   parties,
	})

	responseBody := bytes.NewBuffer(postBody)

	_, err = http.Post(uri, "application/json", responseBody)

	if err != nil {
		//sugar.Error("error creating party with api ", err)
		return err
	}

	return nil
}
