/*
 * GoAstra Backend - Request Validator
 *
 * Centralized input validation using go-playground/validator.
 * Provides custom validation rules and error formatting.
 */
package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

/*
 * Validator wraps go-playground/validator with custom rules.
 */
type Validator struct {
	validate *validator.Validate
}

/*
 * ValidationError represents a single field validation error.
 */
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
}

/*
 * ValidationErrors is a collection of validation errors.
 */
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

/*
 * New creates a new validator instance with custom rules.
 */
func New() *Validator {
	v := validator.New()

	/* Use JSON tag names in error messages */
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	/* Register custom validations */
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("username", validateUsername)

	return &Validator{validate: v}
}

/*
 * Validate validates a struct and returns formatted errors.
 */
func (v *Validator) Validate(i interface{}) *ValidationErrors {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	validationErrors := &ValidationErrors{
		Errors: make([]ValidationError, 0),
	}

	for _, err := range err.(validator.ValidationErrors) {
		validationErrors.Errors = append(validationErrors.Errors, ValidationError{
			Field:   err.Field(),
			Message: formatErrorMessage(err),
			Tag:     err.Tag(),
			Value:   fmt.Sprintf("%v", err.Value()),
		})
	}

	return validationErrors
}

/*
 * ValidateVar validates a single variable.
 */
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

func formatErrorMessage(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "password":
		return fmt.Sprintf("%s does not meet password requirements", field)
	case "username":
		return fmt.Sprintf("%s must contain only letters, numbers, and underscores", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, err.Tag())
	}
}

/*
 * Custom validation functions
 */

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber bool

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		}
	}

	return hasUpper && hasLower && hasNumber
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if len(username) < 3 || len(username) > 32 {
		return false
	}

	for _, char := range username {
		if !isValidUsernameChar(char) {
			return false
		}
	}

	return true
}

func isValidUsernameChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '_'
}

/*
 * Error implements the error interface for ValidationErrors.
 */
func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}

	messages := make([]string, len(ve.Errors))
	for i, err := range ve.Errors {
		messages[i] = err.Message
	}

	return strings.Join(messages, "; ")
}
