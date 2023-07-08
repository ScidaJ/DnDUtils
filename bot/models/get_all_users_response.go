package models

type GetAllUsersResponse struct {
	Message string          `json:"message,omitempty"`
	Data    getAllUsersData `json:"data"`
}

type getAllUsersData struct {
	Data []User `json:"data"`
}
