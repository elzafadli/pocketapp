package custommiddleware

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// TestRecoverHandler tests the RecoverHandler function
func TestRecoverHandler(t *testing.T) {
	// Test cases for different panic values
	tests := []struct {
		name       string
		panicValue interface{}
		expected   string
	}{
		{
			name:       "Should handle error panic",
			panicValue: errors.New("test error"),
			expected:   "test error",
		},
		{
			name:       "Should handle string panic",
			panicValue: "panic message",
			expected:   "panic message",
		},
		{
			name:       "Should handle integer panic",
			panicValue: 123,
			expected:   "123",
		},
		{
			name:       "Should handle nil panic",
			panicValue: nil,
			expected:   "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fiber app
			app := fiber.New()

			// Create a fasthttp request context
			fctx := &fasthttp.RequestCtx{}

			// Set up a basic request
			fctx.Request.SetRequestURI("http://example.com/test")
			fctx.Request.Header.SetMethod("GET")

			// Create a fiber context
			c := app.AcquireCtx(fctx)
			defer app.ReleaseCtx(c)

			// Set a context value to the fiber context
			c.SetUserContext(context.Background())

			// Track the error message passed to golog.Error
			var capturedMsg string
			var capturedErr error

			// Save the original golog.Error function
			originalFunc := gologErrorFunc

			// Replace with a mock function that captures the error
			gologErrorFunc = func(ctx context.Context, msg string, err error, fields ...zap.Field) {
				capturedMsg = msg
				capturedErr = err
			}

			// Restore the original function after the test
			defer func() {
				gologErrorFunc = originalFunc
			}()

			// Call RecoverHandler with the panic value
			RecoverHandler(c, tt.panicValue)

			// Verify that the error message matches what we expect
			assert.Equal(t, tt.expected, capturedMsg, "Error message should match expected value")

			// Verify that the error was properly converted
			if tt.panicValue != nil {
				assert.NotNil(t, capturedErr, "Error should not be nil")
				assert.Equal(t, tt.expected, capturedErr.Error(), "Error string should match expected value")
			}
		})
	}
}

// TestRecoverHandlerWithMiddleware tests using RecoverHandler in a middleware
func TestRecoverHandlerWithMiddleware(t *testing.T) {
	// Save the original golog.Error function
	originalFunc := gologErrorFunc

	// Replace with a mock function that doesn't panic
	gologErrorFunc = func(ctx context.Context, msg string, err error, fields ...zap.Field) {
		// Do nothing in the mock
	}

	// Restore the original function after the test
	defer func() {
		gologErrorFunc = originalFunc
	}()

	// Create a new Fiber app
	app := fiber.New()

	// Add a recover middleware that uses our RecoverHandler
	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				RecoverHandler(c, r)
				_ = c.Status(fiber.StatusInternalServerError).SendString("Recovered from panic")
			}
		}()
		return c.Next()
	})

	// Add a route that will panic
	app.Get("/panic-route", func(c *fiber.Ctx) error {
		panic("test panic")
	})

	// Create a test request
	req := httptest.NewRequest("GET", "/panic-route", nil)

	// Test that the app handles the panic
	resp, err := app.Test(req)
	assert.NoError(t, err, "App.Test should not return an error")
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode, "Response status code should be 500")
}
