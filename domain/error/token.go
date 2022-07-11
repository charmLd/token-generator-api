package error

var (
	UnAuthorizedTokenStr = "unauthorized token"
	ExpiredTokenStr      = "expired token"
	InvalidSignatureStr  = "invalid signature"
)

type UnauthorizedTokenError struct{}

func (e UnauthorizedTokenError) Error() string {
	return UnAuthorizedTokenStr
}

type ExpiredTokenError struct{}

func (e ExpiredTokenError) Error() string {
	return ExpiredTokenStr
}

type InvalidSignature struct{}

func (e InvalidSignature) Error() string {
	return InvalidSignatureStr
}
