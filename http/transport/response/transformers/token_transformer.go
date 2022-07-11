package transformers

import "github.com/charmLd/token-generator-api/domain/entities"

// LoginResponse response for login by email
type TokenResponse struct {
	Status string `json:"status"`
	Token  string `json:"invite_token"`
}

type TokenDetailResponse struct {
	UserId string           `json:"userId"`
	Tokens []entities.Token `json:"tokens"`
}
