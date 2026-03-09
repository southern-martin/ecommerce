package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	pkgerrors "github.com/southern-martin/ecommerce/pkg/errors"
)

// ErrorResponse is the standardized JSON envelope returned for all errors.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the machine-readable code, human-readable message, and
// HTTP status that describes the error.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// RespondError writes a standardized JSON error response derived from err.
// Custom error types from pkg/errors are mapped to their Code and Message
// fields; unknown errors are treated as internal errors.
func RespondError(c *gin.Context, err error) {
	code := "INTERNAL_ERROR"
	message := err.Error()

	switch e := err.(type) {
	case *pkgerrors.NotFoundError:
		code = e.Code
		message = e.Message
	case *pkgerrors.ValidationError:
		code = e.Code
		message = e.Message
	case *pkgerrors.UnauthorizedError:
		code = e.Code
		message = e.Message
	case *pkgerrors.ForbiddenError:
		code = e.Code
		message = e.Message
	case *pkgerrors.ConflictError:
		code = e.Code
		message = e.Message
	case *pkgerrors.InternalError:
		code = e.Code
		message = e.Message
	}

	status := pkgerrors.ToHTTPStatus(err)

	c.JSON(status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Status:  status,
		},
	})
}

// RecoveryWithLogger returns a Gin middleware that recovers from panics,
// logs the incident with a stack trace via zerolog, and returns a 500
// response using the standardized ErrorResponse format.
func RecoveryWithLogger(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()

				logger.Error().
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Str("client_ip", c.ClientIP()).
					Str("user_agent", c.Request.UserAgent()).
					Str("stack", string(stack)).
					Interface("panic", r).
					Msg("recovered from panic")

				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
					Error: ErrorDetail{
						Code:    "INTERNAL_ERROR",
						Message: fmt.Sprintf("internal server error: %v", r),
						Status:  http.StatusInternalServerError,
					},
				})
			}
		}()

		c.Next()
	}
}
