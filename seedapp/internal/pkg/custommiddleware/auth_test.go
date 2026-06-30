package custommiddleware_test

import (
	"io/ioutil"
	"seedapp/config"
	"seedapp/internal/pkg/custommiddleware"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_BasicAuthProtection(t *testing.T) {
	t.Run("should create basic auth middleware with correct users", func(t *testing.T) {
		// Setup
		conf := &config.Config{
			BasicAuths: "user1:pass1,user2:pass2",
		}

		authMiddleware := &custommiddleware.AuthMiddleware{
			Conf: conf,
		}

		// Execute
		handler := authMiddleware.BasicAuthProtection()

		// Assert - we can only verify that a handler is returned
		assert.NotNil(t, handler)
	})

	t.Run("should authenticate with valid credentials", func(t *testing.T) {
		// Setup
		conf := &config.Config{
			BasicAuths: "user1:pass1",
		}

		authMiddleware := &custommiddleware.AuthMiddleware{
			Conf: conf,
		}

		app := fiber.New()
		app.Use(authMiddleware.BasicAuthProtection())
		app.Get("/protected", func(c *fiber.Ctx) error {
			return c.SendString("authenticated")
		})

		// Create a test request with valid basic auth
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Basic dXNlcjE6cGFzczE=") // user1:pass1 in base64

		// Execute
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "authenticated", string(body))
	})

	t.Run("should reject with invalid credentials", func(t *testing.T) {
		// Setup
		conf := &config.Config{
			BasicAuths: "user1:pass1",
		}

		authMiddleware := &custommiddleware.AuthMiddleware{
			Conf: conf,
		}

		app := fiber.New(fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				// In a real application, this would be handled by a proper error handler
				return c.Status(401).SendString("Unauthorized")
			},
		})
		app.Use(authMiddleware.BasicAuthProtection())
		app.Get("/protected", func(c *fiber.Ctx) error {
			return c.SendString("authenticated")
		})

		// Create a test request with invalid basic auth
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Basic aW52YWxpZDppbnZhbGlk") // invalid:invalid in base64

		// Execute
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("should handle multiple users correctly", func(t *testing.T) {
		// Setup
		conf := &config.Config{
			BasicAuths: "user1:pass1,user2:pass2",
		}

		authMiddleware := &custommiddleware.AuthMiddleware{
			Conf: conf,
		}

		app := fiber.New()
		app.Use(authMiddleware.BasicAuthProtection())
		app.Get("/protected", func(c *fiber.Ctx) error {
			return c.SendString("authenticated")
		})

		// Test with user2 credentials
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Basic dXNlcjI6cGFzczI=") // user2:pass2 in base64

		// Execute
		resp, err := app.Test(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "authenticated", string(body))
	})

	t.Run("should handle empty basic auth config", func(t *testing.T) {
		// Setup
		conf := &config.Config{
			BasicAuths: "",
		}

		authMiddleware := &custommiddleware.AuthMiddleware{
			Conf: conf,
		}

		// Execute
		handler := authMiddleware.BasicAuthProtection()

		// Assert - we can only verify that a handler is returned
		assert.NotNil(t, handler)

		// Test that with empty config, all requests are rejected
		app := fiber.New(fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				return c.Status(401).SendString("Unauthorized")
			},
		})
		app.Use(handler)
		app.Get("/protected", func(c *fiber.Ctx) error {
			return c.SendString("authenticated")
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})
}
