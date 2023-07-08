package responses

import "dndutils/bot/models"

type GetAllUsersResponse struct {
	Message string          `json:"message,omitempty"`
	Data    getAllUsersData `json:"data"`
}

type getAllUsersData struct {
	Data []models.User `json:"data"`
}
