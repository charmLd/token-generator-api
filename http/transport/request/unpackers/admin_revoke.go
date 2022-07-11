package unpackers

type AdminRevokeUnpackers struct {
	UserId string `json:"user_id" `
}

func (*AdminRevokeUnpackers) RequiredFormat() string {

	return `
	{
		"user_id":<string>,
	}`
}
