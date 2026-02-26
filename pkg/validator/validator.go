package validator

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	govalidator "github.com/go-playground/validator/v10"

	apperrors "github.com/southern-martin/ecommerce/pkg/errors"
)

var validate *govalidator.Validate

func init() {
	validate = govalidator.New()
}

// ValidateStruct validates a struct using go-playground/validator tags.
// Returns a ValidationError if validation fails.
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(govalidator.ValidationErrors)
	if !ok {
		return apperrors.NewValidationError("VALIDATION_ERROR", err.Error())
	}

	var messages []string
	for _, fe := range validationErrors {
		messages = append(messages, formatFieldError(fe))
	}

	return apperrors.NewValidationError("VALIDATION_ERROR", strings.Join(messages, "; "))
}

// BindAndValidate binds the request body to the given struct and validates it.
func BindAndValidate(c *gin.Context, s interface{}) error {
	if err := c.ShouldBindJSON(s); err != nil {
		return apperrors.NewValidationError("BIND_ERROR", fmt.Sprintf("invalid request body: %s", err.Error()))
	}

	return ValidateStruct(s)
}

// formatFieldError formats a single validation field error into a human-readable message.
func formatFieldError(fe govalidator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("field '%s' is required", fe.Field())
	case "email":
		return fmt.Sprintf("field '%s' must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("field '%s' must be at least %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("field '%s' must be at most %s", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("field '%s' must be exactly %s characters", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("field '%s' must be greater than or equal to %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("field '%s' must be less than or equal to %s", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("field '%s' must be one of [%s]", fe.Field(), fe.Param())
	case "uuid":
		return fmt.Sprintf("field '%s' must be a valid UUID", fe.Field())
	default:
		return fmt.Sprintf("field '%s' failed on '%s' validation", fe.Field(), fe.Tag())
	}
}
