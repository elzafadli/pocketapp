package custommiddleware

import (
	"errors"
	"strings"

	"seedapp/internal/domain/shared/identity"
	"seedapp/internal/pkg/formatter"
	"seedapp/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(codeMap map[error]formatter.Status, statusMap map[error]int) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var httpStatus int
		message := err.Error()
		var errList map[string]interface{}

		// if error is a validator.ErrorMap
		if _err, ok := err.(*validator.ErrorMap); ok {
			message, errList = makeErrorMap(_err.Error())
			err = fiber.ErrBadRequest
		}

		// Retrieve the custom status code if it's a *fiber.Error
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			httpStatus = fiberErr.Code
		} else {
			httpStatus = gethttpstatus(err, statusMap)
		}

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		c.Status(httpStatus)

		code := getcode(err, codeMap)
		traceID, ok := c.Locals("traceId").(string)
		if !ok {
			traceID = identity.NewID().String()
		}

		if len(errList) > 0 {
			return c.JSON(formatter.NewErrorResponseList(code, message, traceID, errList))
		}
		return c.JSON(formatter.NewErrorResponse(code, message, traceID))
	}
}

func makeErrorMap(er string) (string, map[string]interface{}) {
	// Count semicolons to preallocate map with appropriate size
	count := strings.Count(er, ";") + 1
	err := make(map[string]interface{}, count)

	var message string
	errorMsg := strings.Split(er, ";")

	// Use first error message as the main message
	if len(errorMsg) > 0 {
		parts := strings.SplitN(errorMsg[0], ":", 2)
		if len(parts) > 1 {
			message = strings.TrimSpace(parts[1])
			err[strings.TrimSpace(parts[0])] = message
		}
	}

	// Process remaining error messages
	for i := 1; i < len(errorMsg); i++ {
		parts := strings.SplitN(errorMsg[i], ":", 2)
		if len(parts) > 1 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			err[key] = value
		}
	}

	return message, err
}
