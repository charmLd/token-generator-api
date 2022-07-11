package validator

import (
	"context"
	"fmt"
	"github.com/charmLd/token-generator-api/domain/globals"
	"github.com/charmLd/token-generator-api/http/transport/response/transformers"
	"github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	err := validate.RegisterValidation("date-time", IsValidTime)
	if err != nil {
		panic(err)
	}
}

// Validate validates bound values of an unPacker struct against validation rules defined in that unPacker struct.
func Validate(ctx context.Context, data interface{}) []transformers.ErrorTransformer {

	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	// from here you can create your own error messages in whatever language you wish
	errs := err.(validator.ValidationErrors)

	errorArray := make([]transformers.ErrorTransformer, 0)

	for _, e := range errs {
		errorArray = append(errorArray, transformers.ErrorTransformer{
			CorrelationID:    ctx.Value(globals.UUIDKey),
			Code:             "CORE-1000",
			DeveloperMessage: e.Translate(trans),
			Message:          fmt.Sprintf("invalid request param: %v", e.Value()),
		})
	}

	return errorArray
}
