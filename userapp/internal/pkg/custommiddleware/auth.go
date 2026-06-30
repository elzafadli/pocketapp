package custommiddleware

import (
	"errors"
	"strings"

	"userapp/internal/domain/auth"

	"userapp/config"
	"userapp/internal/adapter/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/runsystemid/gocache"
)

type AuthMiddlewareService interface {
	BasicAuthProtection() fiber.Handler
	ApiKeyProtection() fiber.Handler
	JwtProtection() fiber.Handler
	RbacProtection(menuName string, requiredAction string) fiber.Handler
}

type AuthMiddleware struct {
	Conf  *config.Config     `inject:"config"`
	DB    *repository.SqlxDB `inject:"database"`
	Cache gocache.Service    `inject:"cache"`
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

func (j *AuthMiddleware) JwtProtection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		if authHeader == "" {
			return auth.ErrInvalidToken
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return auth.ErrInvalidToken
		}

		tokenString := parts[1]
		secret := j.Conf.JwtSecret
		if secret == "" {
			secret = "default_secret_key"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return auth.ErrInvalidToken
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return auth.ErrInvalidToken
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			return auth.ErrInvalidToken
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			return auth.ErrInvalidToken
		}

		c.Locals("email", email)

		tenantCode, ok := claims["tenant_code"].(string)
		if ok && tenantCode != "" {
			c.Locals("tenant_code", tenantCode)
		}

		roleID, ok := claims["role_id"].(string)
		if ok && roleID != "" {
			c.Locals("role_id", roleID)
		}

		rolesObj, ok := claims["roles"]
		if ok {
			if rolesSlice, ok := rolesObj.([]interface{}); ok {
				roles := make([]string, len(rolesSlice))
				for i, v := range rolesSlice {
					roles[i] = v.(string)
				}
				c.Locals("roles", roles)
			}
		}

		return c.Next()
	}
}

func (j *AuthMiddleware) RbacProtection(menuName string, requiredAction string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleIDObj := c.Locals("role_id")
		if roleIDObj == nil {
			return errors.New("forbidden: insufficient role")
		}

		roleID, ok := roleIDObj.(string)
		if !ok || roleID == "" {
			return errors.New("forbidden: insufficient role")
		}

		var actions []string
		var permissions map[string][]string
		cacheKey := "role-permissions:" + roleID
		err := j.Cache.Get(c.Context(), cacheKey, &permissions)
		if err != nil {
			return errors.New("forbidden: insufficient permissions")
		}
		actions = permissions[menuName]

		for _, action := range actions {
			if action == requiredAction || action == "rbac."+requiredAction {
				return c.Next()
			}
		}

		return errors.New("forbidden: insufficient permissions")
	}
}
