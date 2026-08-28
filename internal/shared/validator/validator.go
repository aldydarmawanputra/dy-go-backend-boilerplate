package validator

import (
	"errors"
	"fmt"

	govalidator "github.com/go-playground/validator/v10"

	"go-backend-boilerplate/internal/shared/apperror"
)

var validate = govalidator.New()

func Struct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	fields := map[string]string{}
	var verrs govalidator.ValidationErrors
	if errors.As(err, &verrs) {
		for _, fe := range verrs {
			fields[fe.Field()] = messageForTag(fe)
		}
	}
	return apperror.Validation("validation failed", fields)
}

func messageForTag(fe govalidator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "url":
		return "must be a valid url"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	default:
		return "invalid value"
	}
}
