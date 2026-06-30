package custommiddleware

import (
	"strings"

	"seedapp/config"
	"seedapp/internal/domain/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

type AuthMiddlewareService interface {
	BasicAuthProtection() fiber.Handler
}

type AuthMiddleware struct {
	Conf *config.Config `inject:"config"`
}

func (j *AuthMiddleware) BasicAuthProtection() fiber.Handler {
	users := strings.Split(j.Conf.BasicAuths, ",")
	basicAuths := make(map[string]string)

	for _, user := range users {
		if user == "" {
			continue
		}

		split := strings.Split(user, ":")
		if len(split) < 2 {
			continue
		}

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
