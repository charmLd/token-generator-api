package unpackers

type ValidateRequest struct {
	Jwt string `json:"jwt" validate:"required"`
}

func (*ValidateRequest) RequiredFormat() string {
	return `
	{
		"status": "SUCCESS”
	}`
}
