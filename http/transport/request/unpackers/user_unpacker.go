package unpackers

type UserUnpackers struct {
	ReferenceId string   `json:"referenceId" validate:"required"`
	App         string   `json:"app" validate:"required"`
	UserType    string   `json:"userType" validate:"required"`
	Roles       []string `json:"roles"`
	DeviceType  string   `json:"deviceType"`
	Name        string   `json:"name" validate:"required"`
	Phone       string   `json:"phone"`
	Password    string   `json:"password" validate:"required"`
	DeviceId    string   `json:"deviceId"`
	Email       string   `json:"email"`
}

// RequiredFormat returns the applicable JSON format for the user_create data structure.
func (*UserUnpackers) RequiredFormat() string {

	return `
	{
      	"ReferenceId":<string>,
      	"deviceType":<string>,
		"name":<string>,
		"userType":<string>,
		"phone":<string>,
      	"appId": <string>,
      	"deviceId":<string>,
      	"email":<string>,
		"password":<string>,
		"roles":[<string>]
    }`
}

type UserUpdateUnpackers struct {
	ReferenceId string `json:"referenceId" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Phone       string `json:"phone"`
	Password    string `json:"password" validate:"required"`
	Email       string `json:"email"`
	DeviceId    string `json:"deviceId"`
	DeviceType  string `json:"deviceType"`
}

// RequiredFormat returns the applicable JSON format for the user_create data structure.
func (*UserUpdateUnpackers) RequiredFormat() string {

	return `
	{
      	"referenceId":<string>,
		"name":<string>,
		"phone":<string>,
      	"email":<string>,
		"deviceId":<string>,
      	"deviceType":<string>,
		"password":<string>
    }`
}

type TokenCreateRequest struct {
	UserID string `json:"user_Id" validate:"required"`
}

// RequiredFormat returns the applicable JSON format for the email login data structure.
func (*TokenCreateRequest) RequiredFormat() string {
	return `
	{
  		"userId": "string",
	}
	`
}
