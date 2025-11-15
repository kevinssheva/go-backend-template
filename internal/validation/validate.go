package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/kevinssheva/go-backend-template/internal/errs"
)

var v = validator.New()

func ValidateStruct(s interface{}) error {
	if err := v.Struct(s); err != nil {
		return errs.New(
			"validation_error",
			400,
			"Invalid request payload",
			errs.WithDetails(extract(err)),
		)
	}

	return nil
}

func extract(err error) map[string]string {
	out := map[string]string{}

	for _, e := range err.(validator.ValidationErrors) {
		out[e.Field()] = e.Error()
	}

	return out
}
