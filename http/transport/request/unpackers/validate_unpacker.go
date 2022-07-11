package unpackers

type ValidateRequest struct {
	UserId      string `json:"user_id" validate:"required"`
	InviteToken string `json:"invite_token" validate:"required"`
}

func (*ValidateRequest) RequiredFormat() string {
	return `
	{
		"user_id": "string",
		"invite_token": "string",
	}`
}
