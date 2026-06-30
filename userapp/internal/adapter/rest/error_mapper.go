package rest

import (
	"userapp/internal/domain/activity_log"
	"userapp/internal/domain/auth"
	"userapp/internal/domain/template"
	"userapp/internal/pkg/formatter"

	"github.com/gofiber/fiber/v2"
)

// CodeMap is a map of error to status code
// ONLY put error that is NOT an internal server error
var CodeMap = map[error]formatter.Status{
	// template
	template.ErrDataNotFound:      formatter.DataNotFound,
	template.ErrDataAlreadyExists: formatter.DataConflict,

	// activity log
	activity_log.ErrDataNotFound:      formatter.DataNotFound,
	activity_log.ErrDataAlreadyExists: formatter.DataConflict,

	//auth
	auth.ErrInvalidBasicAuth: formatter.Unauthorized,
	auth.ErrInvalidApiKey:    formatter.Unauthorized,
	auth.ErrInvalidToken:     formatter.Unauthorized,
}

// StatusMap is a map of error to http status code
// ONLY put error that is NOT an internal server error
var StatusMap = map[error]int{
	// template
	template.ErrDataNotFound:      fiber.StatusNotFound,
	template.ErrDataAlreadyExists: fiber.StatusConflict,

	// activity log
	activity_log.ErrDataNotFound:      fiber.StatusNotFound,
	activity_log.ErrDataAlreadyExists: fiber.StatusConflict,

	//auth
	auth.ErrInvalidBasicAuth: fiber.StatusUnauthorized,
	auth.ErrInvalidApiKey:    fiber.StatusUnauthorized,
	auth.ErrInvalidToken:     fiber.StatusUnauthorized,
}
