package unpackers

type EmailLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RequiredFormat returns the applicable JSON format for the email login data structure.
func (*EmailLoginRequest) RequiredFormat() string {
	return `
	{
  		"email": "merchant@pickme.lk",
		"password": "pass",
	}
	`
}
