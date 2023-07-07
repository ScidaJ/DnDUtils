package models

//TODO: Refactor this. Ugly.

type PostPartyResponse struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    postPartyData
}

type postPartyData struct {
	Data insertedPartyID `json:"data"`
}

type insertedPartyID struct {
	InsertedID string `json:"InsertedID"`
}
