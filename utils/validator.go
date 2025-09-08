package utils

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ValidateMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func mapValidateMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "alphanum":
		return "Should be alphanumeric"
	case "min":
		return "Minimum " + fe.Param()
	case "max":
		return "Maximum " + fe.Param()
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "ltefield":
		return "Should be less than " + fe.Param()
	case "gtefield":
		return "Should be greater than " + fe.Param()
	}

	return "Invalid value"
}

func ValidateError(err error, reflectType reflect.Type, tagName string) []ValidateMessage {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ValidateMessage, len(ve))
		for i, fe := range ve {
			field := fe.Field()
			if structField, ok := reflectType.FieldByName(fe.Field()); ok {
				field = structField.Tag.Get(tagName)
			}
			out[i] = ValidateMessage{field, mapValidateMessage(fe)}
		}
		return out
	}
	return []ValidateMessage{{"", err.Error()}}
}
