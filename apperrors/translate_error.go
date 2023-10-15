package apperrors

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Translate field error to understandable string
func translateError(err validator.FieldError) string {
	var result string
	name := err.Field()
	switch err.Tag() {
	case "required":
		result = name + " is required"
	case "oneof":
		result = name + " should be one of " + err.Param()
	default:
		result = name + " is invalid"
	}
	return result
}

// TryTranslateValidationErrors accepts variable number of errors
// and tries to convert each input error into validator.ValidationErrors.
// If conversion succeeds, it is wrapped inside an instance of AppError.
func TryTranslateValidationErrors(errList ...error) string {
	var allMessages []string
	for _, err := range errList {
		validationErrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return ""
		}

		var messages []string
		for _, e := range validationErrs {
			messages = append(messages, translateError(e))
		}
		allMessages = append(allMessages, messages...)
	}
	message := strings.Join(allMessages[:], ", ")
	return message
}
