package unpackers

type AdminRevokeUnpackers struct {
	InviteToken string `json:"invite_token" `
}

func (*AdminRevokeUnpackers) RequiredFormat() string {

	return `
	{
      	"invite_token":<string>,
	}`
}
