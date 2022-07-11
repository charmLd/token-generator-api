package validator

import (
	//"context"
	//"encoding/base64"
	//"encoding/json"
	//"github.com/charmLd/token-generator-api/domain/entities"
	"github.com/charmLd/token-generator-api/http/error/types"
	validator "gopkg.in/go-playground/validator.v9"
	"net/mail"
	"strconv"
	//"strings"

	"time"
)

// IsValidTime validates a time
func IsValidTime(fl validator.FieldLevel) bool {

	_, err := time.Parse(fl.Param(), fl.Field().String())
	return err == nil
}

func IsEmailValid(email string) error {
	if email != "" {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return (&types.ValidationError{}).New("invalid email")
		}
	}

	return nil
}

func IsUserIDValidUInt(userID string) error {
	_, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return (&types.ValidationError{}).New("user_id must be a valid number")
	}

	return nil
}
func IsBoolValid(boolVal string) error {
	_, err := strconv.ParseBool(boolVal)
	if err != nil {
		return (&types.ValidationError{}).New("invalid boolean value")

	}

	return nil
}
