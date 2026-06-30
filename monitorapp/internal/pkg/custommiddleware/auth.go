package custommiddleware

import (
	"strings"

	"monitorapp/internal/domain/auth"

	"monitorapp/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

type AuthMiddlewareService interface {
	BasicAuthProtection() fiber.Handler
	ApiKeyProtection() fiber.Handler
}

type AuthMiddleware struct {
	Conf *config.Config `inject:"config"`
}

func (j *AuthMiddleware) BasicAuthProtection() fiber.Handler {
	users := strings.Split(j.Conf.BasicAuths, ",")
	basicAuths := make(map[string]string)

	for _, user := range users {
		split := strings.Split(user, ":")
		basicAuths[split[0]] = split[1]
	}

	basicAuth := basicauth.New(basicauth.Config{
		Users:        basicAuths,
		Unauthorized: j.basicAuthUnauthorized,
	})

	return basicAuth
}

func (j *AuthMiddleware) basicAuthUnauthorized(c *fiber.Ctx) error {
	return auth.ErrInvalidBasicAuth
}

func (j *AuthMiddleware) ApiKeyProtection() fiber.Handler {
	return keyauth.New(keyauth.Config{
		KeyLookup: "header:X-Api-Key",
		Validator: j.validateApiKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return auth.ErrInvalidApiKey
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
	})
}

func (j *AuthMiddleware) validateApiKey(c *fiber.Ctx, key string) (bool, error) {
	if strings.TrimSpace(key) == strings.TrimSpace(j.Conf.InternalApiKey) {
		return true, nil
	}

	return false, auth.ErrInvalidApiKey
}
