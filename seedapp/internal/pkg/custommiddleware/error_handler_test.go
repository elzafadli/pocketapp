package custommiddleware

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"seedapp/internal/domain/shared/identity"
	"seedapp/internal/pkg/formatter"
	"seedapp/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	// Define custom errors for testing
	customErr := errors.New("custom error")
	validationErr := validator.NewErrorMap(map[string]error{
		"email":    errors.New("must be a valid email"),
		"password": errors.New("must be at least 8 characters"),
	})

	// Setup code and status maps
	codeMap := map[error]formatter.Status{
		customErr:           formatter.InvalidRequest,
		fiber.ErrBadRequest: formatter.InvalidRequest,
	}
	statusMap := map[error]int{
		customErr:           fiber.StatusBadRequest,
		fiber.ErrBadRequest: fiber.StatusBadRequest,
	}

	// Create a new Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler(codeMap, statusMap),
	})

	// Test case 1: Standard error handling
	app.Get("/standard-error", func(c *fiber.Ctx) error {
		return customErr
	})

	// Test case 2: Fiber error handling
	app.Get("/fiber-error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusForbidden, "forbidden access")
	})

	// Test case 3: Validation error handling
	app.Get("/validation-error", func(c *fiber.Ctx) error {
		return validationErr
	})

	// Test case 4: Unknown error handling
	app.Get("/unknown-error", func(c *fiber.Ctx) error {
		return errors.New("unknown error")
	})

	// Test case 5: With trace ID in context
	app.Get("/with-trace-id", func(c *fiber.Ctx) error {
		traceID := identity.NewID().String()
		c.Locals("traceId", traceID)
		return customErr
	})

	// Run tests
	tests := []struct {
		name           string
		endpoint       string
		expectedStatus int
		expectedCode   string
		checkErrorList bool
	}{
		{
			name:           "Should handle standard error with custom status and code",
			endpoint:       "/standard-error",
			expectedStatus: fiber.StatusBadRequest,
			expectedCode:   formatter.InvalidRequest.String(),
		},
		{
			name:           "Should handle fiber error",
			endpoint:       "/fiber-error",
			expectedStatus: fiber.StatusForbidden,
			expectedCode:   formatter.InternalServerError.String(),
		},
		{
			name:           "Should handle validation error",
			endpoint:       "/validation-error",
			expectedStatus: fiber.StatusBadRequest,
			expectedCode:   formatter.InvalidRequest.String(),
			checkErrorList: true,
		},
		{
			name:           "Should handle unknown error",
			endpoint:       "/unknown-error",
			expectedStatus: fiber.StatusInternalServerError,
			expectedCode:   formatter.InternalServerError.String(),
		},
		{
			name:           "Should include trace ID from context",
			endpoint:       "/with-trace-id",
			expectedStatus: fiber.StatusBadRequest,
			expectedCode:   formatter.InvalidRequest.String(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new http request
			req := httptest.NewRequest("GET", tc.endpoint, nil)
			resp, err := app.Test(req)

			// Assert no error in test execution
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// Parse response body
			var result formatter.Response
			err = json.NewDecoder(resp.Body).Decode(&result)
			assert.NoError(t, err)

			// Assert response fields
			assert.Equal(t, tc.expectedCode, result.Status)
			assert.NotEmpty(t, result.Message)
			assert.NotEmpty(t, result.TraceID)

			// Check for error list if needed
			if tc.checkErrorList {
				assert.NotNil(t, result.ErrorList)
				errorList, ok := result.ErrorList.(map[string]interface{})
				assert.True(t, ok)
				assert.Contains(t, errorList, "email")
				assert.Contains(t, errorList, "password")
			}
		})
	}
}

func TestMakeErrorMap(t *testing.T) {
	tests := []struct {
		name           string
		errorString    string
		expectedMsg    string
		expectedFields []string
	}{
		{
			name:           "Single error",
			errorString:    "email: must be a valid email",
			expectedMsg:    "must be a valid email",
			expectedFields: []string{"email"},
		},
		{
			name:           "Multiple errors",
			errorString:    "email: must be a valid email; password: must be at least 8 characters",
			expectedMsg:    "must be a valid email",
			expectedFields: []string{"email", "password"},
		},
		{
			name:           "Empty error",
			errorString:    "",
			expectedMsg:    "",
			expectedFields: []string{},
		},
		{
			name:           "Malformed error",
			errorString:    "malformed error without colon",
			expectedMsg:    "",
			expectedFields: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			message, errMap := makeErrorMap(tc.errorString)

			// Check main message
			assert.Equal(t, tc.expectedMsg, message)

			// Check error map fields
			for _, field := range tc.expectedFields {
				assert.Contains(t, errMap, field)
			}

			// Check map size
			assert.Equal(t, len(tc.expectedFields), len(errMap))
		})
	}
}
