package validation

import (
	"github.com/go-playground/validator/v10"
)

type validation struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

type errorMap map[string]validation

func newValidation(field string, validationType string) validation {
	return validation{
		Field: field,
		Type:  validationType,
	}
}

func newErrorMap() errorMap {
	return map[string]validation{}
}

func contains(target string, arr []string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}

func Extract(validationErrors validator.ValidationErrors, whiteField []string) errorMap {
	errors := newErrorMap()
	for _, err := range validationErrors {
		var errorMap validation
		field := err.Field()
		if contains(field, whiteField) {
			errorMap = newValidation(field, err.Tag())
		}
		errors[field] = errorMap
	}
	return errors
}
