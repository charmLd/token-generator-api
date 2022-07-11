package transformers

// LoginResponse response for login by email
type LoginResponse struct {
	Status string `json:"status"`
	JWT    string `json:"jwt"`
}
