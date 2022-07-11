package entities

import "time"

const (
	CountryCode = "LK"
)

type Token struct {
	Token          string    `json:"invite_token"`
	UserId         string    `json:"user_id"`
	IsBlacklisted  bool      `json:"is_blacklisted"`
	GeneratedToken string    `json:"jwt_token"`
	CreatedAt      time.Time `json:"created_at"`
	Expiry         time.Time `json:"expiry_at"`
}

type JWTClaims struct {
	UserID   uint64 `json:"user_id"`
	TokenID  string `json:"token_id"`
	UserRole string `json:"user_role"`
	IssuedAt int64  `json:"iat"`
	ExpireAt int64  `json:"exp"`
}

type TokenGenRequest struct {
	UniqueString  string `json:"token"`
	IssuedAt      int64  `json:"iat"`
	ExpireAt      int64  `json:"exp"`
	IsBlacklisted bool   `json:"is_blacklisted"`
}
type TokenDetailsReqParam struct {
	Balcklisted TokenParams
	UserId      string
}

type ValidateRequest struct {
	InviteToken string
	UserId      string
}
type TokenParams struct {
	IsOK  bool
	Value string
}

func (j *TokenGenRequest) Valid() error {
	return nil
}

func (j *JWTClaims) Valid() error {
	return nil
}
